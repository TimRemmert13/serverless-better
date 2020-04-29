package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/serverless/better/lib/db"
	"github.com/serverless/better/lib/model"
)

type Key struct {
	User string    `json:"user"`
	ID   uuid.UUID `json:"id"`
}

type UpdateMapping struct {
	Description string `json:":d"`
	Title       string `json:":t"`
	Achieved    bool   `json:":a"`
}

type PutInput struct {
	Key     Key     `json:"key"`
	Updates Updates `json:"updates"`
}

type Updates struct {
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
func (d *deps) HandleRequest(ctx context.Context, putInput PutInput) (model.Goal, error) {

	// valid input
	empty := uuid.UUID{}
	if putInput.Key.User == "" || putInput.Key.ID == empty {
		return model.Goal{}, model.ResponseError{
			Code:    400,
			Message: "You must provide a valid username and goal id",
		}
	}

	//get  dynamodb session
	if d.ddb == nil {
		d.ddb = db.GetDbSession()
	}

	// map key to attribute values
	key, err := dynamodbattribute.MarshalMap(putInput.Key)

	// map updates
	update, err := dynamodbattribute.MarshalMap(UpdateMapping{
		Description: putInput.Updates.Description,
		Title:       putInput.Updates.Title,
		Achieved:    putInput.Updates.Achieved,
	})

	if err != nil {
		fmt.Println("Problem converting input to attribute map.")
		return model.Goal{}, model.ResponseError{
			Code:    500,
			Message: "Problem converting input to attribute map.",
		}
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
				return model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Exceeded dynamodb provisioned throughput",
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return model.Goal{}, model.ResponseError{
					Code:    404,
					Message: "No goal found with that id.",
				}
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				return model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Goal size is too big to edit and return",
				}
			case dynamodb.ErrCodeRequestLimitExceeded:
				return model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Exceeded dynamodb request limit.",
				}
			default:
				return model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Problem editing goal",
				}
			}
		} else {
			fmt.Println(err.Error())
			return model.Goal{}, model.ResponseError{
				Code:    500,
				Message: "Problem editing goal",
			}
		}
	}

	updates := Updates{}

	// return success
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updates)

	if err != nil {
		fmt.Println(err.Error())
		return model.Goal{}, model.ResponseError{
			Code:    500,
			Message: "Could not create the response.",
		}
	}

	return model.Goal{
		User:        putInput.Key.User,
		ID:          putInput.Key.ID,
		Description: updates.Description,
		Title:       updates.Title,
		Achieved:    updates.Achieved,
	}, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
