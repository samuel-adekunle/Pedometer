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

	"github.com/SamtheSaint/pedometer-api/challenges"
	e "github.com/SamtheSaint/pedometer-api/errors"
	"github.com/SamtheSaint/pedometer-api/responses"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return e.ServerError(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	challengeName := request.PathParameters["challenge_name"]

	resp, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("challenges"),
		Key: map[string]types.AttributeValue{
			"challenge_name": &types.AttributeValueMemberS{
				Value: challengeName,
			},
		},
	})
	if err != nil {
		return e.ServerError(err)
	}

	if resp.Item["challenge_name"] == nil {
		return e.ClientError(404)
	}

	steps := make(map[string]int)
	for user, count := range resp.Item["steps"].(*types.AttributeValueMemberM).Value {
		countVal, err := strconv.ParseInt((count.(*types.AttributeValueMemberN).Value), 10, 0)
		if err != nil {
			return e.ServerError(err)
		}

		steps[user] = int(countVal)
	}

	target, err := strconv.ParseInt(resp.Item["target"].(*types.AttributeValueMemberN).Value, 10, 0)
	if err != nil {
		return e.ServerError(err)
	}

	current, err := strconv.ParseInt(resp.Item["current"].(*types.AttributeValueMemberN).Value, 10, 0)
	if err != nil {
		return e.ServerError(err)
	}

	data := challenges.Challenges{
		ChallengeName: resp.Item["challenge_name"].(*types.AttributeValueMemberS).Value,
		Steps:         steps,
		StartDate:     resp.Item["start_date"].(*types.AttributeValueMemberS).Value,
		EndDate:       resp.Item["end_date"].(*types.AttributeValueMemberS).Value,
		Target:        int(target),
		Current:       int(current),
	}

	jsonData, err := responses.JsonResponseBody(fmt.Sprintf("Retrieved %s from challenges", challengeName), data)
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
