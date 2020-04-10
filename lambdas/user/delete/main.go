package main

import (
	"context"
	"fmt"

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

func HandleDeleteUserRequest(ctx context.Context, deleteUserInput DeleteUserInput) (Response, error) {
	// get cognito service
	svc := cognito.GetCognitoService()

	// create input
	input := &cognitoidentityprovider.DeleteUserInput{
		AccessToken: aws.String(deleteUserInput.AccessToken),
	}

	// attempt to delete the user
	_, err := svc.DeleteUser(input)

	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Unable to delete user."}, err
	}

	return Response{Message: "Successfully deleted user."}, nil
}

func main() {
	lambda.Start(HandleDeleteUserRequest)
}
