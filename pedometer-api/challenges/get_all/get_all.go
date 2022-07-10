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

	"github.com/SamtheSaint/pedometer-api/challenges"
	e "github.com/SamtheSaint/pedometer-api/errors"
	"github.com/SamtheSaint/pedometer-api/responses"
)

func HandleRequest(ctx context.Context) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return e.ServerError(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	resp, err := svc.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String("challenges"),
	})
	if err != nil {
		return e.ServerError(err)
	}

	var data []challenges.Challenges
	for _, v := range resp.Items {
		steps := make(map[string]int)
		for user, count := range v["steps"].(*types.AttributeValueMemberM).Value {
			countVal, err := strconv.ParseInt((count.(*types.AttributeValueMemberN).Value), 10, 0)
			if err != nil {
				return e.ServerError(err)
			}

			steps[user] = int(countVal)
		}

		target, err := strconv.ParseInt(v["target"].(*types.AttributeValueMemberN).Value, 10, 0)
		if err != nil {
			return e.ServerError(err)
		}

		current, err := strconv.ParseInt(v["current"].(*types.AttributeValueMemberN).Value, 10, 0)
		if err != nil {
			return e.ServerError(err)
		}

		data = append(data, challenges.Challenges{
			ChallengeName: v["challenge_name"].(*types.AttributeValueMemberS).Value,
			Steps:         steps,
			StartDate:     v["start_date"].(*types.AttributeValueMemberS).Value,
			EndDate:       v["end_date"].(*types.AttributeValueMemberS).Value,
			Target:        int(target),
			Current:       int(current),
		})
	}

	jsonData, err := responses.JsonResponseBody("Retrieved all from challenges", data)
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
