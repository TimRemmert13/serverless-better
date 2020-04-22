package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type mockLogin struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	Response cognitoidentityprovider.InitiateAuthOutput
}

func (d mockLogin) InitiateAuth(in *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	if *in.AuthParameters["USERNAME"] != "Tim" {
		return nil, awserr.New(
			cognitoidentityprovider.ErrCodeUserNotFoundException,
			"User not found",
			errors.New("User not found"),
		)
	}
	return &d.Response, nil
}

func TestHandleRequest(t *testing.T) {
	t.Run("Successfully authenticate a user", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/login-payload-pos.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var loginInput LoginInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &loginInput)

		// create mock output
		m := mockLogin{
			Response: cognitoidentityprovider.InitiateAuthOutput{
				AuthenticationResult: &cognitoidentityprovider.AuthenticationResultType{
					AccessToken:  aws.String("2342543fervwrfvwbrthwt"),
					ExpiresIn:    aws.Int64(1000),
					IdToken:      aws.String("id token"),
					RefreshToken: aws.String("refresh token"),
				},
			},
		}

		// create dependancy object
		d := deps{
			cognito: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, loginInput)

		if err != nil {
			t.Error("Error when trying to successfully login")
		}

		if *result.AccessToken == "" ||
			*result.ExpiresIn == 0 ||
			*result.IDToken == "" ||
			*result.RefreshToken == "" {
			t.Error("Sent invalid successful login response")
		}
	})

	t.Run("Catch user not found exception", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/login-payload-neg.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var loginInput LoginInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &loginInput)

		// create mock output
		m := mockLogin{}

		// create dependancy object
		d := deps{
			cognito: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, loginInput)

		if result.Message != "Username or password is incorrect" || err == nil {
			t.Error("Unable to catch and handle user not found exception")
		}
	})
}
