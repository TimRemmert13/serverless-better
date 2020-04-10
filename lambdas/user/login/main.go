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

func HandleLoginRequest(ctx context.Context, loginInput LoginInput) (Response, error) {
	// get cognito service
	svc := cognito.GetCognitoService()
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

	output, err := svc.InitiateAuth(input)

	// process results
	if err != nil {
		fmt.Println(err.Error())
		return Response{Message: "Problem authenticating user"}, err
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
	//HandleLoginRequest(nil, LoginInput{Username: "Tim", Password: "Morty2012!"})
	lambda.Start(HandleLoginRequest)
}
