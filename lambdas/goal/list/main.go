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
	Message string       `json:"result"`
	Goals   []model.Goal `json:"goals"`
}

type QueryMapping struct {
	User string `json:":u"`
}

type ListInput struct {
	User string `json:"user"`
}

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and add it to dynamodb
*/
func HandleRequest(ctx context.Context, listInput ListInput) (Response, error) {

	// get local dynamodb session
	db := db.GetLocalSession()

	// map updates
	query, err := dynamodbattribute.MarshalMap(QueryMapping{
		User: listInput.User,
	})

	if err != nil {
		fmt.Println("Problem converting input to attribute map.")
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String("Goals"),
		ExpressionAttributeNames:  aws.StringMap(map[string]string{"#U": "user"}),
		KeyConditionExpression:    aws.String("#U = :u"),
		ExpressionAttributeValues: query,
	}

	result, err := db.Query(input)

	// handle possible errors
	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Problem executing query."}, err
	}

	// return success
	goals := []model.Goal{}
	dynamodbattribute.UnmarshalListOfMaps(result.Items, &goals)

	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Could not create the response."}, nil
	}

	response := Response{
		Message: "Successfully retrieved goal!",
		Goals:   goals,
	}
	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
