package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/serverless/better/lib/cognito"
	"github.com/serverless/better/lib/util"
)

type UserInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Response struct {
	Message string `json:"response"`
}

func handleRequest(ctx context.Context, userInput UserInput) (Response, error) {
	svc := cognito.GetCognitoService()

	email := "email"
	name := "name"

	emailAttribute := &cognitoidentityprovider.AttributeType{
		Name:  &email,
		Value: &userInput.Email,
	}

	nameAttribute := &cognitoidentityprovider.AttributeType{
		Name:  &name,
		Value: &userInput.Name,
	}

	// configure create user input
	input := &cognitoidentityprovider.SignUpInput{
		ClientId:       aws.String("16btm3mund4aaeupitgjljmual"),
		Password:       aws.String(userInput.Password),
		SecretHash:     aws.String(util.GenerateSecretHash(userInput.Name)),
		UserAttributes: []*cognitoidentityprovider.AttributeType{emailAttribute, nameAttribute},
		Username:       aws.String(userInput.Name),
	}

	output, err := svc.SignUp(input)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(output)
	return Response{Message: "successfully created user"}, nil
}

func main() {
	lambda.Start(handleRequest)
}
