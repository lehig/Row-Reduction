//go:build lambda

package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// LambdaHandler handles the Lambda function invocation
func LambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle CORS preflight
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type",
			},
			Body: "",
		}, nil
	}

	// Only allow POST requests
	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
				"Content-Type":                "application/json",
			},
			Body: `{"error":"Method not allowed"}`,
		}, nil
	}

	// Parse request body
	var req RREFRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
				"Content-Type":                "application/json",
			},
			Body: fmt.Sprintf(`{"error":"Invalid request: %v"}`, err),
		}, nil
	}

	// Calculate RREF
	rref := RREF(req.Matrix.Data)

	response := RREFResponse{
		Original: req.Matrix,
		RREF:     Matrix{Data: rref},
	}

	// Marshal response
	responseBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin": "*",
				"Content-Type":                "application/json",
			},
			Body: fmt.Sprintf(`{"error":"Error encoding response: %v"}`, err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Content-Type":                "application/json",
		},
		Body: string(responseBody),
	}, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
