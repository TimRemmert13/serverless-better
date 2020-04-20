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
	"github.com/serverless/better/lib/model"
)

type Response struct {
	Message     string `json:"result"`
	DeletedGoal model.Goal
}

type Key struct {
	User string `json:"user"`
	ID   string `json:"id"`
}

type deps struct {
	ddb dynamodbiface.DynamoDBAPI
}

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and add it to dynamodb
*/
func (d *deps) HandleRequest(ctx context.Context, inputKey Key) (Response, error) {

	// validate input
	if inputKey.ID == "" || inputKey.User == "" {
		return Response{}, errors.New("You must provide a goal id and username")
	}

	if d.ddb == nil {
		d.ddb = db.GetDbSession()
	}

	// map key to attribute values
	key, err := dynamodbattribute.MarshalMap(inputKey)

	input := &dynamodb.DeleteItemInput{
		Key:          key,
		TableName:    aws.String("Goals"),
		ReturnValues: aws.String("ALL_OLD"),
	}
	result, err := d.ddb.DeleteItem(input)

	// handle exceptions
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return Response{"Exceeded provisioned throughput for request", model.Goal{}}, err
			case dynamodb.ErrCodeResourceNotFoundException:
				return Response{"No goal found by that id", model.Goal{}}, err
			case dynamodb.ErrCodeRequestLimitExceeded:
				return Response{"Dynamodb request limit has been reached", model.Goal{}}, err
			default:
				return Response{"Problem deleting item", model.Goal{}}, err
			}
		}
	}

	// return success
	deletedGoal := model.Goal{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &deletedGoal)

	if err != nil {
		fmt.Println(err)
		return Response{Message: "Could not create the response."}, err
	}

	response := Response{
		Message:     "Successfully deleted the goal",
		DeletedGoal: deletedGoal,
	}
	return response, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
