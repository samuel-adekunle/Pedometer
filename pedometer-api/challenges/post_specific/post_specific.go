package main

import (
	"context"
	"encoding/json"
	"fmt"

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

// HandleRequest will create a new challenge
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return e.ServerError(err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	challengeName := request.PathParameters["challenge_name"]

	var jsonBody challenges.ChallengeRequestBody
	err = json.Unmarshal([]byte(request.Body), &jsonBody)
	if err != nil {
		return e.ClientError(400) // Malformed json
	}

	if challengeName != jsonBody.ChallengeName {
		return e.ClientError(403) // Unauthorised request
	}

	resp, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
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

	if resp.Item["challenge_name"] != nil {
		return e.ClientError(403) // challenge already exists
	}

	_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("challenges"),
		Item: map[string]types.AttributeValue{
			"challenge_name": &types.AttributeValueMemberS{
				Value: challengeName,
			},
			"current": &types.AttributeValueMemberN{
				Value: "0",
			},
			"target": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%v", jsonBody.Target),
			},
			"start_date": &types.AttributeValueMemberS{
				Value: jsonBody.StartDate,
			},
			"end_date": &types.AttributeValueMemberS{
				Value: jsonBody.EndDate,
			},
			"steps": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{},
			},
		},
	})
	if err != nil {
		return e.ServerError(err)
	}

	jsonData, err := responses.JsonResponseBody(fmt.Sprintf("Updated %s in challenges", challengeName), nil)

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
