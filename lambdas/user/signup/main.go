package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/serverless/better/lib/cognito"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
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

type deps struct {
	cognito cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

func (d *deps) HandleRequest(ctx context.Context, userInput UserInput) (Response, error) {

	// validate input
	if userInput.Name == "" || userInput.Email == "" || userInput.Password == "" {
		return Response{}, errors.New("You must provide a username, email, and password")
	}

	// initialize cognito service
	if d.cognito == nil {
		d.cognito = cognito.GetCognitoService()
	}

	emailAttribute := &cognitoidentityprovider.AttributeType{
		Name:  aws.String("email"),
		Value: &userInput.Email,
	}

	nameAttribute := &cognitoidentityprovider.AttributeType{
		Name:  aws.String("name"),
		Value: &userInput.Name,
	}

	// configure create user input
	input := &cognitoidentityprovider.SignUpInput{
		ClientId:       aws.String(os.Getenv("AWS_CLIENT_ID")),
		Password:       aws.String(userInput.Password),
		SecretHash:     aws.String(util.GenerateSecretHash(userInput.Name)),
		UserAttributes: []*cognitoidentityprovider.AttributeType{emailAttribute, nameAttribute},
		Username:       aws.String(userInput.Name),
	}

	_, err := d.cognito.SignUp(input)

	// handle possible exceptions
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeInvalidParameterException:
				return Response{Message: "Invalid password"}, err
			case cognitoidentityprovider.ErrCodeUsernameExistsException:
				return Response{Message: "User already exists"}, err
			case cognitoidentityprovider.ErrCodeTooManyRequestsException:
				return Response{Message: "Too many request made to cognito for user signup"}, err
			default:
				fmt.Println(aerr.Error())
				return Response{Message: "Problem signing up the user"}, err
			}
		} else {
			fmt.Println(err.Error())
			return Response{Message: "Problem signing up the user"}, err
		}
	}

	return Response{Message: "successfully created user"}, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
