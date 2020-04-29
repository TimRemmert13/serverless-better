package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/serverless/better/lib/db"
	"github.com/serverless/better/lib/model"
)

type GetInput struct {
	User string    `json:"user"`
	ID   uuid.UUID `json:"id"`
}

type deps struct {
	ddb dynamodbiface.DynamoDBAPI
}

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and add it to dynamodb
*/
func (d *deps) HandleRequest(ctx context.Context, getInput GetInput) (model.Goal, error) {

	// valid input
	empty := uuid.UUID{}
	if getInput.User == "" || getInput.ID == empty {
		return model.Goal{}, model.ResponseError{
			Code:    400,
			Message: "Request must have a user and an id",
		}
	}

	// map key to attribute values
	key, err := dynamodbattribute.MarshalMap(getInput)

	if err != nil {
		fmt.Println("Problem converting input to attribute map.")
		return model.Goal{}, model.ResponseError{
			Code:    500,
			Message: "Problem converting input to attribute map",
		}
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
				return model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "dynamodb provisioned throughput exceeded",
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return model.Goal{}, model.ResponseError{
					Code:    404,
					Message: "No goal found by that id",
				}
			case dynamodb.ErrCodeRequestLimitExceeded:
				return model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Exceeded dynamodb request limit",
				}
			default:
				return model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Problem retrieving the goal",
				}
			}
		} else {
			fmt.Println(err.Error())
			return model.Goal{}, model.ResponseError{
				Code:    500,
				Message: "Problem retrieving the goal",
			}
		}
	}

	// return success
	goal := model.Goal{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &goal)

	if err != nil {
		fmt.Println(err.Error())
		return model.Goal{}, model.ResponseError{
			Code:    500,
			Message: "Could not create the response.",
		}
	}

	return goal, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
