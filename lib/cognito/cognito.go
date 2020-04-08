package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/aws/aws-sdk-go/aws/session"
)

func GetCognitoService() *cognitoidentityprovider.CognitoIdentityProvider {
	cogSession := session.Must(session.NewSession())

	// Create a CognitoIdentityProvider client from session
	svc := cognitoidentityprovider.New(cogSession, aws.NewConfig().WithRegion("us-east-1"))

	return svc
}
