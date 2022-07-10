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
	challengeName := request.PathParameters["challenge_name"]

	resp, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("steps"),
		Key: map[string]types.AttributeValue{
			"user_name": &types.AttributeValueMemberS{
				Value: userName,
			},
		},
		ProjectionExpression: aws.String("user_name"),
	})
	if err != nil {
		return e.ServerError(err)
	}

	if resp.Item["user_name"] == nil {
		return e.ClientError(404) //User does not exist
	}

	resp, err = svc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("challenges"),
		Key: map[string]types.AttributeValue{
			"challenge_name": &types.AttributeValueMemberS{
				Value: challengeName,
			},
		},
		ProjectionExpression: aws.String("challenge_name"),
	})
	if err != nil {
		return e.ServerError(err)
	}

	if resp.Item["challenge_name"] == nil {
		return e.ClientError(404) //Challenge does not exist
	}

	_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("steps"),
		Key: map[string]types.AttributeValue{
			"user_name": &types.AttributeValueMemberS{
				Value: userName,
			},
		},
		UpdateExpression: aws.String("DELETE challenges :c"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":c": &types.AttributeValueMemberSS{
				Value: []string{
					challengeName,
				},
			},
		},
		ReturnValues: types.ReturnValueNone,
	})
	if err != nil {
		return e.ServerError(err)
	}

	jsonData, err := responses.JsonResponseBody(fmt.Sprintf("Deleted %s challenge from %s in steps", challengeName, userName), nil)

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
