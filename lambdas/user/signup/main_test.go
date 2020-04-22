package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

type mockSignup struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	Response cognitoidentityprovider.SignUpOutput
}

func (d mockSignup) Signup(in *cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, error) {
	if *in.Username != "Tim" {
		return nil, awserr.New(
			cognitoidentityprovider.ErrCodeUsernameExistsException,
			"user already exist",
			errors.New("User already exists"),
		)
	}
	return &d.Response, nil
}

func TestHandleRequest(t *testing.T) {
	t.Run("Successfully signup new user", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/signup-payload.json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var userInput UserInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &userInput)
		err = jsonFile.Close()
		if err != nil {
			fmt.Println(err)
		}

		// create mock output
		m := mockSignup{
			Response: cognitoidentityprovider.SignUpOutput{
				UserConfirmed: aws.Bool(false),
				UserSub:       aws.String("uuid"),
			},
		}

		// create dependancy object
		d := deps{
			cognito: m,
		}

		//execute test of function
		_, err = d.HandleRequest(nil, userInput)

		if err != nil {
			t.Error("Failed to signup new user")
		}
	})

	t.Run("User already exists exception", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/signup-payload.json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer jsonFile.Close()
		var userInput UserInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &userInput)
		userInput.Name = "AlreadyExists"

		// create mock output
		m := mockSignup{}

		// create dependancy object
		d := deps{
			cognito: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, userInput)

		if result.Message != "User already exists" || err == nil {
			t.Error("Failed to catch and handle user already exists exception")
		}
	})
}
