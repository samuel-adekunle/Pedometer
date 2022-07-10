package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	e "github.com/SamtheSaint/pedometer-api/errors"
	"github.com/SamtheSaint/pedometer-api/responses"
	"github.com/SamtheSaint/pedometer-api/steps"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return e.ServerError(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	userName := request.PathParameters["user_name"]

	resp, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
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

	if resp.Item["user_name"] == nil {
		return e.ClientError(404) // resource does not exist
	}

	var challenges []string
	for _, c := range resp.Item["challenges"].(*types.AttributeValueMemberSS).Value {
		challenges = append(challenges, c)
	}

	count := make(map[string]int)
	for date, distance := range resp.Item["count"].(*types.AttributeValueMemberM).Value {
		distanceInt64, err := strconv.ParseInt((distance.(*types.AttributeValueMemberN).Value), 10, 0)
		if err != nil {
			return e.ServerError(err)
		}

		count[date] = int(distanceInt64)
	}

	target, err := strconv.ParseInt(resp.Item["daily_target"].(*types.AttributeValueMemberN).Value, 10, 0)
	if err != nil {
		return e.ServerError(err)
	}

	data := steps.Step{
		UserName:    resp.Item["user_name"].(*types.AttributeValueMemberS).Value,
		Challenges:  challenges,
		Count:       count,
		DailyTarget: int(target),
	}

	jsonData, err := responses.JsonResponseBody(fmt.Sprintf("Retrieved %s from steps", userName), data)
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
