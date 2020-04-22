package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/serverless/better/lib/cognito"
)

type SignOutInput struct {
	AccessToken string `json:"access_token"`
}

type Response struct {
	Message string `json:"result"`
}

type deps struct {
	cognito cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

func (d *deps) HandleRequest(ctx context.Context, signOutInput SignOutInput) (Response, error) {

	// validate input
	if signOutInput.AccessToken == "" {
		return Response{}, errors.New("You must provide a valid access token")
	}
	// get cognito session
	if d.cognito == nil {
		d.cognito = cognito.GetCognitoService()
	}

	// create sign out input
	input := &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(signOutInput.AccessToken),
	}

	// initiate sign out
	_, err := d.cognito.GlobalSignOut(input)

	// handle possible exceptions
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeResourceNotFoundException:
				return Response{Message: "The access token provided is invalid"}, err
			case cognitoidentityprovider.ErrCodeTooManyRequestsException:
				return Response{Message: "Too many request made to validate the code"}, err
			default:
				fmt.Println(aerr.Error())
				return Response{Message: "Problem signing out user"}, err
			}
		} else {
			fmt.Println(err.Error())
			return Response{Message: "Problem signing out user"}, err
		}
	}

	return Response{Message: "Successfully signed out user."}, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
