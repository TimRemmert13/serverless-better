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

type mockGetItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.GetItemOutput
}

func (d mockGetItem) GetItem(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
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
	t.Run("Successfully Get Goal", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/get-payload-pos.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var getInput GetInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &getInput)

		retrievedGoal := model.Goal{
			User:        "Tim",
			ID:          "1",
			Title:       "my new title",
			Created:     time.Now(),
			Description: "my new description",
			Achieved:    true,
		}

		av, _ := dynamodbattribute.MarshalMap(retrievedGoal)

		// create mock output
		m := mockGetItem{
			Response: dynamodb.GetItemOutput{
				Item: av,
			},
		}

		// create dependancy object
		d := deps{
			ddb: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, getInput)

		if err != nil {
			t.Error("geting a goal failed")
		}

		if result.Message != "Successfully retrieved goal!" {
			t.Error("Incorrect Successful message sent")
		}
	})

	t.Run("Goal not Found", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/get-payload-neg.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var getInput GetInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &getInput)

		// create mock output
		m := mockGetItem{}

		// create dependancy object
		d := deps{
			ddb: m,
		}

		// execute test of function
		result, err := d.HandleRequest(nil, getInput)

		if result.Message != "No goal found by that id" || err == nil {
			t.Error("Not catching and handling a not found exception")
		}
	})
}
