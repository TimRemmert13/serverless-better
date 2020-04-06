package main

import (
	"context"
	"fmt"

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

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and add it to dynamodb
*/
func HandleRequest(ctx context.Context, inputKey Key) (Response, error) {

	// map key to attribute values
	key, err := dynamodbattribute.MarshalMap(inputKey)

	// // get local dynamodb session
	db := db.GetLocalSession()

	input := &dynamodb.DeleteItemInput{
		Key:          key,
		TableName:    aws.String("Goals"),
		ReturnValues: aws.String("ALL_OLD"),
	}

	result, err := db.DeleteItem(input)

	// handle possible errors
	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Problem saving changes."}, err
	}

	// return success
	deletedGoal := model.Goal{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &deletedGoal)

	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Could not create the response."}, nil
	}

	response := Response{
		Message:     "Successfully deleted the goal",
		DeletedGoal: deletedGoal,
	}
	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
