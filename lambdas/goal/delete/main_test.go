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

type mockDeleteItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.DeleteItemOutput
}

func (d mockDeleteItem) DeleteItem(in *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	if *in.Key["id"].S != "1" {
		return nil, awserr.New(
			dynamodb.ErrCodeResourceNotFoundException,
			"resource not found",
			errors.New("Resource not found error"),
		)
	}
	return &d.Response, nil
}

func TestHandleRequest(t *testing.T) {
	t.Run("Successful Deletion", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/delete-payload-pos.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var key Key
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &key)

		deletedGoal := model.Goal{
			User:        "Tim",
			ID:          "1",
			Created:     time.Now(),
			Description: "my description",
			Achieved:    false,
		}

		av, _ := dynamodbattribute.MarshalMap(deletedGoal)

		// create mock output
		m := mockDeleteItem{
			Response: dynamodb.DeleteItemOutput{
				Attributes: av,
			},
		}

		// create dependancy object
		d := deps{
			ddb: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, key)

		if err != nil {
			t.Error("Deleting a goal failed")
		}

		if result.Message != "Successfully deleted the goal" {
			t.Error("Incorrect Successfully message sent")
		}
	})

	t.Run("Goal not found", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/delete-payload-neg.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var key Key
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &key)

		// create mock output
		m := mockDeleteItem{}

		// create dependancy object
		d := deps{
			ddb: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, key)

		if result.Message != "No goal found by that id" || err == nil {
			t.Error("Failed to catch and handle resource not found exception")
		}
	})
}
