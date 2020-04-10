package main

import (
	"context"
	"fmt"

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

func HandleLogoutRequest(ctx context.Context, signOutInput SignOutInput) (Response, error) {
	// get cognito session
	svc := cognito.GetCognitoService()

	// create sign out input
	input := &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(signOutInput.AccessToken),
	}

	// initiate sign out
	_, err := svc.GlobalSignOut(input)

	// process results
	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Unable to signout user."}, err
	}

	return Response{Message: "Successfully signed out user."}, nil
}

func main() {
	lambda.Start(HandleLogoutRequest)
}
