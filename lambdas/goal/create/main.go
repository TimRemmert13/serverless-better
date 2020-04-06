package main

import (
	"context"
	"fmt"

	"github.com/serverless/better/lib/model"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/serverless/better/lib/db"
)

type Response struct {
	Message string `json:"response"`
}

// type PostInput struct {
// 	User        string `json:"user"`
// 	ID          string `json:"id"`
// 	Description string `json:"description"`
// 	Title       string `json:"title"`
// 	Achieved    bool   `json:"achieved"`
// 	Created     string `json:"created"`
// 	Updated     string `json:"updated"`
// }

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and add it to dynamodb
*/
func HandleRequest(ctx context.Context, goal model.Goal) (Response, error) {

	// map input to function to a goal struct
	// goal := model.Goal{
	// 	User:        postInput.User,
	// 	ID:          postInput.ID,
	// 	Description: postInput.Description,
	// 	Title:       postInput.Title,
	// 	Achieved:    postInput.Achieved,
	// 	Created:     postInput.Created,
	// 	Updated:     postInput.Updated,
	// }

	// // get local dynamodb session
	db := db.GetLocalSession()

	// // map go struct to dynamodb attribute values
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
	result, err := db.PutItem(input)

	// handle possible errors
	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Problem saving changes."}, err
	}

	// return success
	fmt.Println(result)
	return Response{Message: "Successfully added new goal!"}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
