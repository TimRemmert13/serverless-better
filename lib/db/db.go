package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
}

func GetLocalSession() *dynamodb.DynamoDB {
	region := "us-east-1"
	endpoint := "http://DynamoDBEndpoint:8000"
	sess, _ := session.NewSession(&aws.Config{
		Region:   &region,
		Endpoint: &endpoint,
	})
	return dynamodb.New(sess)
}
