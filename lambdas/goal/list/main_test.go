package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/serverless/better/lib/model"
)

type mockQuery struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.QueryOutput
}

func (d mockQuery) Query(in *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if *in.ExpressionAttributeValues["user"].S != "Tim" {
		return nil, awserr.New(
			dynamodb.ErrCodeResourceNotFoundException,
			"resource not found",
			errors.New("Resource not found error"),
		)
	}
	return &d.Response, nil
}

func TestHandleRequest(t *testing.T) {
	t.Run("Successfully Get all users goals", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/list-payload-pos.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var queryInput ListInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &queryInput)

		retrievedGoals := []model.Goal{
			model.Goal{
				User:        "Tim",
				ID:          "1",
				Title:       "my title",
				Created:     time.Now(),
				Description: "my description",
				Achieved:    true,
			},
			model.Goal{
				User:        "Tim",
				ID:          "2",
				Title:       "my new title",
				Created:     time.Now(),
				Description: "my new description",
				Achieved:    false,
			},
		}

		avmap := []map[string]*dynamodb.AttributeValue{}
		for i := range retrievedGoals {
			avmap[i], _ = dynamodbattribute.MarshalMap(retrievedGoals[i])
		}

		// create mock output
		m := mockQuery{
			Response: dynamodb.QueryOutput{
				Items: avmap,
			},
		}

		// create dependancy object
		d := deps{
			ddb: m,
		}

		//execute test of function
		_, err = d.HandleRequest(nil, queryInput)

		if err != nil {
			t.Error("geting all users goals failed")
		}
	})
}
