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
	Message string     `json:"result"`
	Goal    model.Goal `json:"goal"`
}

type GetInput struct {
	User string `json:"user"`
	ID   string `json:"id"`
}

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and add it to dynamodb
*/
func HandleRequest(ctx context.Context, getInput GetInput) (Response, error) {

	// map key to attribute values
	key, err := dynamodbattribute.MarshalMap(getInput)

	if err != nil {
		fmt.Println("Problem converting input to attribute map.")
		fmt.Println(err.Error())
	}

	// get local dynamodb session
	db := db.GetLocalSession()

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String("Goals"),
	}

	result, err := db.GetItem(input)

	// handle possible errors
	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Problem saving changes."}, err
	}

	// return success
	goal := model.Goal{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &goal)

	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Could not create the response."}, nil
	}

	response := Response{
		Message: "Successfully retrieved goal!",
		Goal:    goal,
	}
	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
