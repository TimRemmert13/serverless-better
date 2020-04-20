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
	Message string     `json:"result"`
	Goal    model.Goal `json:"goal"`
}

type GetInput struct {
	User string `json:"user"`
	ID   string `json:"id"`
}

type deps struct {
	ddb dynamodbiface.DynamoDBAPI
}

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and add it to dynamodb
*/
func (d *deps) HandleRequest(ctx context.Context, getInput GetInput) (Response, error) {

	// valid input
	if getInput.User == "" || getInput.ID == "" {
		return Response{}, errors.New("Request must have a user and id")
	}

	// map key to attribute values
	key, err := dynamodbattribute.MarshalMap(getInput)

	if err != nil {
		fmt.Println("Problem converting input to attribute map.")
		return Response{}, errors.New("Problem converting input to attribute map")
	}

	// get local dynamodb session
	if d.ddb == nil {
		d.ddb = db.GetDbSession()
	}

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String("Goals"),
	}

	result, err := d.ddb.GetItem(input)

	// handle possible errors
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return Response{"dynamodb provisioned throughput exceeded", model.Goal{}}, err
			case dynamodb.ErrCodeResourceNotFoundException:
				return Response{"No goal found by that id", model.Goal{}}, err
			case dynamodb.ErrCodeRequestLimitExceeded:
				return Response{"Exceeded dynamodb request limit", model.Goal{}}, err
			default:
				return Response{"Problem retrieving the goal", model.Goal{}}, err
			}
		} else {
			fmt.Println(err.Error())
			return Response{"Problem retrieving the goal", model.Goal{}}, err
		}
	}

	// return success
	goal := model.Goal{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &goal)

	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Could not create the response."}, err
	}

	response := Response{
		Message: "Successfully retrieved goal!",
		Goal:    goal,
	}
	return response, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
