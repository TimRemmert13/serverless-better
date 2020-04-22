package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type mockLogout struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	Response cognitoidentityprovider.GlobalSignOutOutput
}

func (d mockLogout) GlobalSignOut(in *cognitoidentityprovider.GlobalSignOutInput) (*cognitoidentityprovider.GlobalSignOutOutput, error) {
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
		jsonFile, err := os.Open("./testdata/logout-payload.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var signOutInput SignOutInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &signOutInput)

		// create mock output
		m := mockLogout{}

		// create dependancy object
		d := deps{
			cognito: m,
		}

		//execute test of function
		_, err = d.HandleRequest(nil, signOutInput)

		if err != nil {
			t.Error("Failed to logout user")
		}
	})

	t.Run("Send incorrect token for logout", func(t *testing.T) {

		// load test data
		signOutInput := SignOutInput{AccessToken: "IncorrectToken"}

		// create mock output
		m := mockLogout{}

		// create dependancy object
		d := deps{
			cognito: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, signOutInput)

		if result.Message != "The access token provided is invalid" || err == nil {
			t.Error("Failed to catch and handle and incorrect token exception")
		}
	})
}
