package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/serverless/better/lib/cognito"
	"github.com/serverless/better/lib/util"
)

type ConfirmInput struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type Response struct {
	Message string `json:"result"`
}

func HandleRequest(ctx context.Context, confirmInput ConfirmInput) (Response, error) {
	// get cognito service
	svc := cognito.GetCognitoService()

	// create verify input
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(os.Getenv("AWS_CLIENT_ID")),
		ConfirmationCode: aws.String(confirmInput.Token),
		SecretHash:       aws.String(util.GenerateSecretHash(confirmInput.Username)),
		Username:         aws.String(confirmInput.Username),
	}

	// call verify token
	_, err := svc.ConfirmSignUp(input)

	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Problme verifying the token"}, err
	}

	// process results
	return Response{Message: "Successfully verified the token"}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
