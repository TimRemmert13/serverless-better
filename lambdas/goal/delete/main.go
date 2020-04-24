package main

import (
	"context"
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
func (d *deps) HandleRequest(ctx context.Context, inputKey Key) (model.Goal, error) {

	// validate input
	if inputKey.ID == "" || inputKey.User == "" {
		return model.Goal{}, model.ResponseError{
			Code:    400,
			Message: "You must provide a valid username and goal id",
		}
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
				return model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Exceeded provisioned throughput for request",
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return model.Goal{}, model.ResponseError{
					Code:    404,
					Message: "No goal found by that id",
				}
			case dynamodb.ErrCodeRequestLimitExceeded:
				return model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Dynamodb request limit has been reached",
				}
			default:
				return model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Problem deleting item",
				}
			}
		}
	}

	// return success
	deletedGoal := model.Goal{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &deletedGoal)

	if err != nil {
		fmt.Println(err)
		return model.Goal{}, model.ResponseError{
			Code:    500,
			Message: "Could not create the response.",
		}
	}

	return deletedGoal, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
