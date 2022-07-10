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
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return e.ServerError(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	challengeName := request.PathParameters["challenge_name"]
	userName := request.PathParameters["user_name"]

	resp, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("challenges"),
		Key: map[string]types.AttributeValue{
			"challenge_name": &types.AttributeValueMemberS{
				Value: challengeName,
			},
		},
		ProjectionExpression: aws.String("challenge_name, steps"),
	})
	if err != nil {
		return e.ServerError(err)
	}

	if resp.Item["challenge_name"] == nil {
		return e.ClientError(404)
	}

	userSteps := resp.Item["steps"].(*types.AttributeValueMemberM).Value[userName]
	if userSteps == nil {
		return e.ClientError(404)
	}

	count, err := strconv.ParseInt(userSteps.(*types.AttributeValueMemberN).Value, 10, 0)
	if err != nil {
		return e.ServerError(err)
	}

	data := struct {
		UserName      string `json:"user_name"`
		ChallengeName string `json:"challenge_name"`
		Steps         int    `json:"steps"`
	}{
		UserName:      userName,
		ChallengeName: challengeName,
		Steps:         int(count),
	}

	jsonData, err := responses.JsonResponseBody(fmt.Sprintf("Retrieved %s steps from %s challenge", userName, challengeName), data)
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
