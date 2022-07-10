package main

import (
	"context"
	"encoding/json"
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

// HandleRequest will create a new user and update values
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return e.ServerError(err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	userName := request.PathParameters["user_name"]

	var jsonBody steps.StepRequestBody
	err = json.Unmarshal([]byte(request.Body), &jsonBody)
	if err != nil {
		return e.ClientError(400) // Malformed json
	}

	if userName != jsonBody.UserName {
		return e.ClientError(403) // Unauthorised request
	}

	var totalCount int
	for _, c := range jsonBody.Count {
		totalCount += c
	}

	if jsonBody.DailyTarget == 0 {
		jsonBody.DailyTarget = 6000 // default daily target value, 0 causes division errors later
	}

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

	newUser := resp.Item["user_name"] == nil

	if newUser {
		count := make(map[string]types.AttributeValue)
		for date, distance := range jsonBody.Count {
			count[date] = &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%v", distance),
			}
		}

		item := map[string]types.AttributeValue{
			"user_name": &types.AttributeValueMemberS{
				Value: userName,
			},
			"challenges": &types.AttributeValueMemberSS{
				Value: []string{
					"default",
				},
			},
			"count": &types.AttributeValueMemberM{
				Value: count,
			},
			"daily_target": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%v", jsonBody.DailyTarget),
			},
		}

		_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
			TableName:    aws.String("steps"),
			Item:         item,
			ReturnValues: types.ReturnValueNone,
		})
		if err != nil {
			return e.ServerError(err)
		}

		_, err = svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String("challenges"),
			Key: map[string]types.AttributeValue{
				"challenge_name": &types.AttributeValueMemberS{
					Value: "default",
				},
			},
			UpdateExpression: aws.String("SET steps.#u = :c, #w = #w + :c"),
			ExpressionAttributeNames: map[string]string{
				"#u": userName,
				"#w": "current",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":c": &types.AttributeValueMemberN{
					Value: fmt.Sprintf("%v", totalCount),
				},
			},
		})
		if err != nil {
			return e.ServerError(err)
		}

		jsonData, err := responses.JsonResponseBody(fmt.Sprintf("Updated %s in steps", userName), nil)
		if err != nil {
			return e.ServerError(err)
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(jsonData),
			Headers:    responses.JsonResponseHeader(),
		}, nil
	}

	for _, c := range resp.Item["challenges"].(*types.AttributeValueMemberSS).Value {
		_, err := svc.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String("challenges"),
			Key: map[string]types.AttributeValue{
				"challenge_name": &types.AttributeValueMemberS{
					Value: c,
				},
			},
			UpdateExpression: aws.String("SET steps.#u = steps.#u + :c, #w = #w + :c"),
			ExpressionAttributeNames: map[string]string{
				"#u": userName,
				"#w": "current",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":c": &types.AttributeValueMemberN{
					Value: fmt.Sprintf("%v", totalCount),
				},
			},
		})
		if err != nil {
			return e.ServerError(err)
		}
	}

	var prevCountVal int64
	for date, amount := range jsonBody.Count {

		prevCount := resp.Item["count"].(*types.AttributeValueMemberM).Value[date]
		if prevCount == nil {
			prevCountVal = 0
		} else {
			prevCountVal, err = strconv.ParseInt(prevCount.(*types.AttributeValueMemberN).Value, 10, 0)
			if err != nil {
				return e.ServerError(err)
			}
		}

		resp.Item["count"].(*types.AttributeValueMemberM).Value[date] = &types.AttributeValueMemberN{
			Value: fmt.Sprint(amount + int(prevCountVal)),
		}
	}

	resp.Item["daily_target"] = &types.AttributeValueMemberN{
		Value: fmt.Sprintf("%v", jsonBody.DailyTarget),
	}

	_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:    aws.String("steps"),
		Item:         resp.Item,
		ReturnValues: types.ReturnValueNone,
	})
	if err != nil {
		return e.ServerError(err)
	}

	jsonData, err := responses.JsonResponseBody(fmt.Sprintf("Updated %s in steps", userName), nil)
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
