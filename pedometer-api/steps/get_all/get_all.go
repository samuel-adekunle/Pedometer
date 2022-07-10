package main

import (
	"context"
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

func HandleRequest(ctx context.Context) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return e.ServerError(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	resp, err := svc.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String("steps"),
	})
	if err != nil {
		return e.ServerError(err)
	}

	var data []steps.Step
	for _, item := range resp.Items {
		var challenges []string
		for _, c := range item["challenges"].(*types.AttributeValueMemberSS).Value {
			challenges = append(challenges, c)
		}

		count := make(map[string]int)
		for date, distance := range item["count"].(*types.AttributeValueMemberM).Value {
			distanceInt64, err := strconv.ParseInt((distance.(*types.AttributeValueMemberN).Value), 10, 0)
			if err != nil {
				return e.ServerError(err)
			}

			count[date] = int(distanceInt64)
		}

		target, err := strconv.ParseInt(item["daily_target"].(*types.AttributeValueMemberN).Value, 10, 0)
		if err != nil {
			return e.ServerError(err)
		}

		data = append(data, steps.Step{
			UserName:    item["user_name"].(*types.AttributeValueMemberS).Value,
			Challenges:  challenges,
			Count:       count,
			DailyTarget: int(target),
		})
	}

	jsonData, err := responses.JsonResponseBody("Retrieved all from steps", data)
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
