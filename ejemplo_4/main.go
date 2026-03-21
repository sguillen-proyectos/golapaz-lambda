package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Method: %s Path: %s\n", request.HTTPMethod, request.Path)

	switch {
	case request.HTTPMethod == "GET" && request.Path == "/hello":
		return handleGet(request)
	case request.HTTPMethod == "POST" && request.Path == "/hello":
		return handlePost(request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"message":"route not found"}`,
		}, nil
	}
}

func handleGet(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := request.QueryStringParameters["name"]
	if name == "" {
		name = "World"
	}

	body := fmt.Sprintf(`{"message":"Hello, %s!"}`, name)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}, nil
}

func handlePost(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var payload map[string]string
	if err := json.Unmarshal([]byte(request.Body), &payload); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message":"invalid JSON body"}`,
		}, nil
	}

	name := payload["name"]
	if name == "" {
		name = "World"
	}

	body := fmt.Sprintf(`{"message":"Hello, %s! (POST)"}`, name)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}, nil
}

func main() {
	runLocally := os.Getenv("RUN_LOCALLY")

	if runLocally != "" {
		fmt.Println("Running locally - no Lambda runtime detected")
		return
	}
	lambda.Start(Handler)
}
