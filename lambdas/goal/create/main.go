package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

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
	if goal.User == "" || goal.ID == "" || goal.Title == "" {
		return Response{}, errors.New("Missing property user, id, or title in goal")
	}

	goal.Created = time.Now()

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
				return Response{Message: "Encounted provisioned throughput exception for dynamodb"}, err
			case dynamodb.ErrCodeResourceNotFoundException:
				return Response{Message: "Could not find the specified table"}, err
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				return Response{Message: "Size of the goal is too large"}, err
			case dynamodb.ErrCodeRequestLimitExceeded:
				return Response{Message: "Reached the request limit for dynamodb"}, err
			default:
				return Response{Message: "Problem creating a new goal"}, err
			}
		} else {
			fmt.Println(err.Error())
			return Response{Message: "Problem creating a new goal"}, err
		}
	}

	// return success and log results in cloud watch
	return Response{Message: "Successfully added new goal!"}, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
