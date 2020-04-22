package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type mockDelete struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	Response cognitoidentityprovider.DeleteUserOutput
}

func (d mockDelete) DeleteUser(in *cognitoidentityprovider.DeleteUserInput) (*cognitoidentityprovider.DeleteUserOutput, error) {
	if *in.AccessToken != "CorrectToken" {
		return nil, awserr.New(
			cognitoidentityprovider.ErrCodeResourceNotFoundException,
			"Resources not found",
			errors.New("Resources not found"),
		)
	}
	return &d.Response, nil
}

func TestHandleRequest(t *testing.T) {
	t.Run("Successfully logout a user", func(t *testing.T) {

		// load test data
		deleteUserInput := DeleteUserInput{AccessToken: "CorrectToken"}

		// create mock output
		m := mockDelete{}

		// create dependancy object
		d := deps{
			cognito: m,
		}

		//execute test of function
		_, err := d.HandleRequest(nil, deleteUserInput)

		if err != nil {
			t.Error("Failed to delete user")
		}
	})

	t.Run("delete user attempt with invalid token", func(t *testing.T) {

		// load test data
		deleteUserInput := DeleteUserInput{AccessToken: "IncorrectToken"}

		// create mock output
		m := mockDelete{}

		// create dependancy object
		d := deps{
			cognito: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, deleteUserInput)

		if result.Message != "Invalid access token provided" || err == nil {
			t.Error("Failed to catch and handle invalid token exception")
		}
	})
}
