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

type mockUpdateItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.UpdateItemOutput
}

func (d mockUpdateItem) UpdateItem(in *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
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
	t.Run("Successful Update", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/update-payload-pos.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var putInput PutInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &putInput)

		updatedGoal := model.Goal{
			User:        "Tim",
			ID:          "1",
			Title:       "my new title",
			Created:     time.Now(),
			Description: "my new description",
			Achieved:    true,
		}

		av, _ := dynamodbattribute.MarshalMap(updatedGoal)

		// create mock output
		m := mockUpdateItem{
			Response: dynamodb.UpdateItemOutput{
				Attributes: av,
			},
		}

		// create dependancy object
		d := deps{
			ddb: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, putInput)

		if err != nil {
			t.Error("Updating a goal failed")
		}

		if result.Message != "Successfully updated the following fields" {
			t.Error("Incorrect Successful message sent")
		}
	})

	t.Run("Cant find goal to update", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/update-payload-neg.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var putInput PutInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &putInput)

		// create mock output
		m := mockUpdateItem{}

		// create dependancy object
		d := deps{
			ddb: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, putInput)

		if result.Message != "No goal found with that id." || err == nil {
			t.Error("Failed to catch and handle resource not found exception")
		}
	})
}
