package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GetDbSession() *dynamodb.DynamoDB {

	// if running lambdas locally for developement use dynamodb instance
	// in docker container
	if os.Getenv("RUN_ENV") == "LOCAL" {
		region := "us-east-1"
		endpoint := "http://localhost:8000"
		sess, _ := session.NewSession(&aws.Config{
			Region:   &region,
			Endpoint: &endpoint,
		})
		return dynamodb.New(sess)
	}

	// else return remote production session of dynamodb
	region := "us-east-1"
	sess, _ := session.NewSession(&aws.Config{
		Region: &region,
	})
	return dynamodb.New(sess)
}
