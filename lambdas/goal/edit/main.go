package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/serverless/better/lib/db"
)

type Response struct {
	Message     string `json:"result"`
	UpdatedGoal UpdatedGoal
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

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and add it to dynamodb
*/
func HandleRequest(ctx context.Context, putInput PutInput) (Response, error) {

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
	}

	// // get local dynamodb session
	db := db.GetLocalSession()

	input := &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String("Goals"),
		UpdateExpression:          aws.String("set description = :d, title = :t, achieved = :a"),
		ExpressionAttributeValues: update,
		ReturnValues:              aws.String("UPDATED_NEW"),
	}

	result, err := db.UpdateItem(input)

	// handle possible errors
	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Problem saving changes."}, err
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
	lambda.Start(HandleRequest)
}
