package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
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
	}

	// create dynamodb input
	input := &dynamodb.PutItemInput{
		Item:                   av,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String("Goals"),
	}

	//attempt insert into dynamodb
	result, err := d.ddb.PutItem(input)

	// handle possible errors
	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Problem saving changes."}, err
	}

	// return success and log results in cloud watch
	fmt.Println(result)
	return Response{Message: "Successfully added new goal!"}, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
