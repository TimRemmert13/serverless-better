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

type DeleteUserInput struct {
	AccessToken string `json:"access_token"`
}

type Response struct {
	Message string `json:"result"`
}

type deps struct {
	cognito cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

func (d *deps) HandleRequest(ctx context.Context, deleteUserInput DeleteUserInput) (Response, error) {
	// validate input
	if deleteUserInput.AccessToken == "" {
		return Response{}, errors.New("You must provide a valid access token")
	}
	// get cognito service
	if d.cognito == nil {
		d.cognito = cognito.GetCognitoService()
	}

	// create input
	input := &cognitoidentityprovider.DeleteUserInput{
		AccessToken: aws.String(deleteUserInput.AccessToken),
	}

	// attempt to delete the user
	_, err := d.cognito.DeleteUser(input)

	// handle any exceptions
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeResourceNotFoundException:
				return Response{Message: "Invalid access token provided"}, err
			case cognitoidentityprovider.ErrCodePasswordResetRequiredException:
				return Response{Message: "You must reset your password before you can delete your account"}, err
			case cognitoidentityprovider.ErrCodeUserNotConfirmedException:
				return Response{Message: "You must first confirm your email before you can delete your account"}, err
			case cognitoidentityprovider.ErrCodeTooManyRequestsException:
				return Response{Message: "Too many request made to login"}, err
			case cognitoidentityprovider.ErrCodeUserNotFoundException:
				return Response{Message: "No user found"}, err
			default:
				fmt.Println(aerr.Error())
				return Response{Message: "Problem deleting user"}, err
			}
		} else {
			fmt.Println(err.Error())
			return Response{Message: "Problem deleting user"}, err
		}
	}

	return Response{Message: "Successfully deleted user."}, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
