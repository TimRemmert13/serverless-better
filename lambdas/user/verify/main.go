package main

import (
	"context"
	"fmt"
	"os"

	"github.com/serverless/better/lib/model"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
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

type deps struct {
	cognito cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

func (d *deps) HandleRequest(ctx context.Context, confirmInput ConfirmInput) (Response, error) {

	// validate input
	if confirmInput.Username == "" || confirmInput.Token == "" {
		return Response{}, model.ResponseError{
			Code:    400,
			Message: "You must provide a username and token",
		}
	}

	// get cognito service
	if d.cognito == nil {
		d.cognito = cognito.GetCognitoService()
	}

	// create verify input
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(os.Getenv("AWS_CLIENT_ID")),
		ConfirmationCode: aws.String(confirmInput.Token),
		SecretHash:       aws.String(util.GenerateSecretHash(confirmInput.Username)),
		Username:         aws.String(confirmInput.Username),
	}

	// call verify token
	_, err := d.cognito.ConfirmSignUp(input)

	// handle possible exceptions
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeTooManyFailedAttemptsException:
				return Response{}, model.ResponseError{
					Code:    500,
					Message: "Too many failed attempts to verify the code",
				}
			case cognitoidentityprovider.ErrCodeCodeMismatchException:
				return Response{}, model.ResponseError{
					Code:    400,
					Message: "Provided incorrect code",
				}
			case cognitoidentityprovider.ErrCodeExpiredCodeException:
				return Response{}, model.ResponseError{
					Code:    400,
					Message: "The code has expired",
				}
			case cognitoidentityprovider.ErrCodeTooManyRequestsException:
				return Response{}, model.ResponseError{
					Code:    500,
					Message: "Too many request made to validate the code",
				}
			case cognitoidentityprovider.ErrCodeUserNotFoundException:
				return Response{}, model.ResponseError{
					Code:    404,
					Message: "No user found",
				}
			default:
				fmt.Println(aerr.Error())
				return Response{}, model.ResponseError{
					Code:    500,
					Message: "Problem verifying code",
				}
			}
		} else {
			fmt.Println(err.Error())
			return Response{}, model.ResponseError{
				Code:    500,
				Message: "Problem verifying code",
			}
		}
	}

	// process results
	return Response{Message: "Successfully verified the code"}, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
