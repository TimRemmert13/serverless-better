package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/serverless/better/lib/db"
)

type Response struct {
	Message     string      `json:"result"`
	UpdatedGoal UpdatedGoal `json:"updated"`
}

type UpdatedGoal struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	Acheived    bool   `json:"achieved"`
}

type Key struct {
	User string `json:"user"`
	ID   string `json:"id"`
}

type UpdateMapping struct {
	Description string `json:":d"`
	Title       string `json:":t"`
	Achieved    bool   `json:":a"`
}

type PutInput struct {
	Key         Key    `json:"key"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Achieved    bool   `json:"achieved"`
}

type deps struct {
	ddb dynamodbiface.DynamoDBAPI
}

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and add it to dynamodb
*/
func (d *deps) HandleRequest(ctx context.Context, putInput PutInput) (Response, error) {

	// valid input
	if putInput.Key.User == "" || putInput.Key.ID == "" {
		return Response{}, errors.New("You must provide a valid username and id")
	}

	//get  dynamodb session
	if d.ddb == nil {
		d.ddb = db.GetDbSession()
	}

	// map key to attribute values
	key, err := dynamodbattribute.MarshalMap(putInput.Key)

	// map updates
	update, err := dynamodbattribute.MarshalMap(UpdateMapping{
		Description: putInput.Description,
		Title:       putInput.Title,
		Achieved:    putInput.Achieved,
	})

	if err != nil {
		fmt.Println("Problem converting input to attribute map.")
		return Response{}, err
	}

	input := &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String("Goals"),
		UpdateExpression:          aws.String("set description = :d, title = :t, achieved = :a"),
		ExpressionAttributeValues: update,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	result, err := d.ddb.UpdateItem(input)

	// handle possible errors
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return Response{Message: "Provisioned throughput has been exceeded for dynamodb"}, err
			case dynamodb.ErrCodeResourceNotFoundException:
				return Response{Message: "No goal found with that id."}, err
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				return Response{Message: "Goal size is too big to edit and return"}, err
			case dynamodb.ErrCodeRequestLimitExceeded:
				return Response{Message: "Exceeded dynamodb request limit."}, err
			default:
				return Response{Message: "Problem editing goal"}, err
			}
		} else {
			fmt.Println(err.Error())
			return Response{Message: "Problem editing goal"}, err
		}
	}

	// return success
	updatedGoal := UpdatedGoal{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedGoal)

	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Could not create the response."}, nil
	}

	response := Response{
		Message:     "Successfully updated the following fields",
		UpdatedGoal: updatedGoal,
	}
	return response, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
