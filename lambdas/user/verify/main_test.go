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

type mockVerify struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI
	Response cognitoidentityprovider.ConfirmSignUpOutput
}

func (d mockVerify) ConfirmSignUp(in *cognitoidentityprovider.ConfirmSignUpInput) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	if *in.ConfirmationCode != "183609" {
		return nil, awserr.New(
			cognitoidentityprovider.ErrCodeExpiredCodeException,
			"code has expired",
			errors.New("Code has expired"),
		)
	}
	return &d.Response, nil
}

func TestHandleRequest(t *testing.T) {
	t.Run("Successfully verify granted code", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/verify-payload.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var confirmInput ConfirmInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &confirmInput)

		// create mock output
		m := mockVerify{
			Response: cognitoidentityprovider.ConfirmSignUpOutput{},
		}

		// create dependancy object
		d := deps{
			cognito: m,
		}

		//execute test of function
		_, err = d.HandleRequest(nil, confirmInput)

		if err != nil {
			t.Error("verifying confirm code failed")
		}
	})

	t.Run("Generated code expired", func(t *testing.T) {

		// load test data
		jsonFile, err := os.Open("./testdata/verify-payload.json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		var confirmInput ConfirmInput
		byteJSON, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteJSON, &confirmInput)
		confirmInput.Token = "12345"

		// create mock output
		m := mockVerify{
			Response: cognitoidentityprovider.ConfirmSignUpOutput{},
		}

		// create dependancy object
		d := deps{
			cognito: m,
		}

		//execute test of function
		result, err := d.HandleRequest(nil, confirmInput)

		if result.Message != "The code has expired" || err == nil {
			t.Error("failed to catch and handle expired code exception")
		}
	})
}
