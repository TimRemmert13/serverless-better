package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

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

type deps struct {
	ddb dynamodbiface.DynamoDBAPI
}

/* HandleRequest is a function for lambda function to take an input of a goal in a json form
and add it to dynamodb
*/
func (d *deps) HandleRequest(ctx context.Context, listInput ListInput) (Response, error) {

	// validate input
	if listInput.User == "" {
		return Response{}, errors.New("You must include a user in your request")
	}

	// get local dynamodb session
	if d.ddb == nil {
		d.ddb = db.GetDbSession()
	}

	// create query mapping
	query, err := dynamodbattribute.MarshalMap(QueryMapping{
		User: listInput.User,
	})

	if err != nil {
		fmt.Println(err)
		return Response{"Problem creating query attribute map", []model.Goal{}}, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String("Goals"),
		ExpressionAttributeNames:  aws.StringMap(map[string]string{"#U": "user"}),
		KeyConditionExpression:    aws.String("#U = :u"),
		ExpressionAttributeValues: query,
	}

	result, err := d.ddb.Query(input)

	// handle possible errors
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return Response{Message: "Dyanamodb provisioned throughput limit reached", Goals: []model.Goal{}}, err
			case dynamodb.ErrCodeResourceNotFoundException:
				return Response{Message: "User not found", Goals: []model.Goal{}}, err
			case dynamodb.ErrCodeRequestLimitExceeded:
				return Response{Message: "Reached dynamodb request limit", Goals: []model.Goal{}}, err
			default:
				fmt.Println(aerr.Error())
				return Response{Message: "Problem getting all goals for the user", Goals: []model.Goal{}}, err
			}
		} else {
			fmt.Println(err.Error())
			return Response{Message: "Problem getting all goals for the user", Goals: []model.Goal{}}, err
		}
	}

	// return success
	goals := []model.Goal{}
	dynamodbattribute.UnmarshalListOfMaps(result.Items, &goals)

	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Problem unmarshalling response from dynamodb", Goals: []model.Goal{}}, err
	}

	response := Response{
		Goals: goals,
	}
	return response, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
