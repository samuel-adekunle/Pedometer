package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	e "github.com/SamtheSaint/pedometer-api/errors"
	"github.com/SamtheSaint/pedometer-api/responses"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return e.ServerError(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	userName := request.PathParameters["user_name"]

	// Note: operation is idempotent so doesn't matter if key does not exists
	_, err = svc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String("steps"),
		Key: map[string]types.AttributeValue{
			"user_name": &types.AttributeValueMemberS{
				Value: userName,
			},
		},
	})
	if err != nil {
		return e.ServerError(err)
	}

	jsonData, err := responses.JsonResponseBody(fmt.Sprintf("Deleted %s from steps", userName), nil)
	if err != nil {
		return e.ServerError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
		Headers:    responses.JsonResponseHeader(),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
