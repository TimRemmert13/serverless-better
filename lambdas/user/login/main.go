package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/serverless/better/lib/cognito"
	"github.com/serverless/better/lib/util"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Message      string  `json:"result"`
	AccessToken  *string `json:"access_token"`
	ExpiresIn    *int64  `json:"expires"`
	IDToken      *string `json:"id_token"`
	RefreshToken *string `json:"refresh_token"`
}

type deps struct {
	cognito cognitoidentityprovideriface.CognitoIdentityProviderAPI
}

func (d *deps) HandleRequest(ctx context.Context, loginInput LoginInput) (Response, error) {
	// validate input
	if loginInput.Username == "" || loginInput.Password == "" {
		return Response{}, errors.New("You must provide a valid username and password")
	}

	// get cognito service
	if d.cognito == nil {
		d.cognito = cognito.GetCognitoService()
	}

	// create auth input
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String(os.Getenv("AWS_CLIENT_ID")),
		AuthParameters: aws.StringMap(map[string]string{
			"SECRET_HASH": util.GenerateSecretHash(loginInput.Username),
			"USERNAME":    loginInput.Username,
			"PASSWORD":    loginInput.Password,
		}),
	}

	output, err := d.cognito.InitiateAuth(input)

	// handle possible exceptions
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodePasswordResetRequiredException:
				return Response{Message: "You must reset your password before you can login"}, err
			case cognitoidentityprovider.ErrCodeUserNotConfirmedException:
				return Response{Message: "You have not confirmed your email address yet"}, err
			case cognitoidentityprovider.ErrCodeInvalidUserPoolConfigurationException:
				return Response{Message: "Cognito userpool not configured for this request"}, err
			case cognitoidentityprovider.ErrCodeTooManyRequestsException:
				return Response{Message: "Too many request made to login"}, err
			case cognitoidentityprovider.ErrCodeUserNotFoundException:
				return Response{Message: "Username or password is incorrect"}, err
			default:
				fmt.Println(aerr.Error())
				return Response{Message: "Problem authenticating user"}, err
			}
		} else {
			fmt.Println(err.Error())
			return Response{Message: "Problem authenticating user"}, err
		}
	}

	response := Response{
		Message:      "Successfully Authenticated user.",
		AccessToken:  output.AuthenticationResult.AccessToken,
		ExpiresIn:    output.AuthenticationResult.ExpiresIn,
		IDToken:      output.AuthenticationResult.IdToken,
		RefreshToken: output.AuthenticationResult.RefreshToken,
	}
	return response, nil
}

func main() {
	d := deps{}
	lambda.Start(d.HandleRequest)
}
