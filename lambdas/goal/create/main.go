package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"

	"github.com/serverless/better/lib/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/serverless/better/lib/db"
)

type Response struct {
	Message string `json:"response"`
}

type deps struct {
	ddb dynamodbiface.DynamoDBAPI
}

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and adds it to dynamodb for storage
*/
func (d *deps) HandleRequest(ctx context.Context, goal model.Goal) (Response, error) {

	// validate input
	if goal.User == "" || goal.Title == "" {
		return Response{}, model.ResponseError{
			Code:    400,
			Message: "You must provide a goal user, id, and title",
		}
	}

	goal.Created = time.Now()
	goal.ID = uuid.New()

	// get dynamodb session
	if d.ddb == nil {
		d.ddb = db.GetDbSession()
	}

	// map go struct to dynamodb attribute values
	av, err := dynamodbattribute.MarshalMap(goal)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// create dynamodb input
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Goals"),
	}

	//attempt insert into dynamodb
	_, err = d.ddb.PutItem(input)

	// handle possible errors
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return Response{}, model.ResponseError{
					Code:    400,
					Message: "You must provide a goal user, id, and title",
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return Response{}, model.ResponseError{
					Code:    404,
					Message: "Could not find the specified table",
				}
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				return Response{}, model.ResponseError{
					Code:    500,
					Message: "Goal size is too large to update",
				}
			case dynamodb.ErrCodeRequestLimitExceeded:
				return Response{}, model.ResponseError{
					Code:    500,
					Message: "Reached the request limit for dynamodb",
				}
			default:
				fmt.Println(aerr)
				return Response{}, model.ResponseError{
					Code:    500,
					Message: "Problem creating a new goal",
				}
			}
		} else {
			fmt.Println(err.Error())
			return Response{}, model.ResponseError{
				Code:    500,
				Message: "Problem creating a new goal",
			}
		}
	}

	// return success and log results in cloud watch
	return Response{Message: "Successfully added new goal!"}, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
