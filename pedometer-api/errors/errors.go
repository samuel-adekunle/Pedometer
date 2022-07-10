package errors

import (
	"log"
	"net/http"

	"github.com/SamtheSaint/pedometer-api/responses"
	"github.com/aws/aws-lambda-go/events"
)

func ServerError(err error) (events.APIGatewayProxyResponse, error) {
	log.Println(err) // print server error to cloudwatch

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
		Headers:    responses.JsonResponseHeader(),
	}, nil
}

func ClientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
		Headers:    responses.JsonResponseHeader(),
	}, nil
}
