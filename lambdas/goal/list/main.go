package main

import (
	"context"
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
func (d *deps) HandleRequest(ctx context.Context, listInput ListInput) ([]model.Goal, error) {

	// validate input
	if listInput.User == "" {
		return []model.Goal{}, model.ResponseError{
			Code:    400,
			Message: "You must include a user in your request",
		}
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
		return []model.Goal{}, model.ResponseError{
			Code:    500,
			Message: "Problem creating query attribute map",
		}
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
				return []model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Dyanamodb provisioned throughput limit reached",
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return []model.Goal{}, model.ResponseError{
					Code:    404,
					Message: "User not found",
				}
			case dynamodb.ErrCodeRequestLimitExceeded:
				return []model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Reached dynamodb request limit",
				}
			default:
				fmt.Println(aerr.Error())
				return []model.Goal{}, model.ResponseError{
					Code:    500,
					Message: "Problem getting all goals for the user",
				}
			}
		} else {
			fmt.Println(err.Error())
			return []model.Goal{}, model.ResponseError{
				Code:    500,
				Message: "Problem getting all goals for the user",
			}
		}
	}

	// return success
	goals := []model.Goal{}
	dynamodbattribute.UnmarshalListOfMaps(result.Items, &goals)

	if err != nil {
		fmt.Println(err.Error())
		return []model.Goal{}, model.ResponseError{
			Code:    500,
			Message: "Problem unmarshalling response from dynamodb",
		}
	}

	return goals, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
