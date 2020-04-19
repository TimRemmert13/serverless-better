package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/serverless/better/lib/model"
)

type mockPutItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.PutItemOutput
}

func (d mockPutItem) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &d.Response, nil
}

func TestHandleRequest(t *testing.T) {
	t.Run("Successful Request", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/create-payload-pos.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var goal model.Goal
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &goal)

		// create mock output
		m := mockPutItem{
			Response: dynamodb.PutItemOutput{
				ConsumedCapacity: &dynamodb.ConsumedCapacity{
					CapacityUnits: aws.Float64(2),
					TableName:     aws.String("Goals"),
				},
			},
		}

		// create dependancy object
		d := deps{
			ddb: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, goal)

		if err != nil {
			t.Error("Creating a goal successfully failed")
		}

		if result.Message != "Successfully added new goal!" {
			t.Error("Incorrect Successfully added new goal message sent")
		}

	})

	// negative scneario invalid input
	t.Run("Invalid input Request", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/create-payload-neg.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var goal model.Goal
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &goal)

		// create mock output
		m := mockPutItem{
			Response: dynamodb.PutItemOutput{
				ConsumedCapacity: &dynamodb.ConsumedCapacity{
					CapacityUnits: aws.Float64(2),
					TableName:     aws.String("Goals"),
				},
			},
		}

		// create dependency object
		d := deps{
			ddb: m,
		}

		//execute test of function
		_, err = d.HandleRequest(nil, goal)

		if err.Error() != "Missing property user, id, or title in goal" {
			t.Error("Invalid input passed to create input not caught")
		}

	})
}
