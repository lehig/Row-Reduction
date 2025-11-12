//go:build lambda

package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestRREF(t *testing.T) {
	tests := []struct {
		name     string
		input    [3][3]float64
		expected [3][3]float64
	}{
		{
			name: "Identity matrix",
			input: [3][3]float64{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			expected: [3][3]float64{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
		},
		{
			name: "Simple 2x2 system",
			input: [3][3]float64{
				{2, 1, 5},
				{1, 1, 3},
				{0, 0, 0},
			},
			expected: [3][3]float64{
				{1, 0, 2},
				{0, 1, 1},
				{0, 0, 0},
			},
		},
		{
			name: "Zero matrix",
			input: [3][3]float64{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			expected: [3][3]float64{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
		},
		{
			name: "Singular matrix (linearly dependent rows)",
			input: [3][3]float64{
				{1, 2, 3},
				{2, 4, 6},
				{3, 6, 9},
			},
			expected: [3][3]float64{
				{1, 2, 3},
				{0, 0, 0},
				{0, 0, 0},
			},
		},
		{
			name: "Matrix with negative numbers",
			input: [3][3]float64{
				{1, -1, 0},
				{2, 1, 3},
				{-1, 2, 1},
			},
			expected: [3][3]float64{
				{1, 0, 1},
				{0, 1, 1},
				{0, 0, 0},
			},
		},
		{
			name: "Matrix with decimals",
			input: [3][3]float64{
				{1.5, 2.5, 3.5},
				{3, 5, 7},
				{4.5, 7.5, 10.5},
			},
			// This matrix has linearly dependent rows (row 3 = 3 * row 1, row 2 = 2 * row 1)
			// So RREF will have zeros in rows 2 and 3
			expected: [3][3]float64{
				{1, 1.6666666666666667, 2.3333333333333335},
				{0, 0, 0},
				{0, 0, 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RREF(tt.input)
			if !matricesEqual(result, tt.expected) {
				t.Errorf("RREF() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function to compare matrices with floating point tolerance
func matricesEqual(a, b [3][3]float64) bool {
	const epsilon = 1e-9
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			diff := a[i][j] - b[i][j]
			if diff < 0 {
				diff = -diff
			}
			if diff > epsilon {
				return false
			}
		}
	}
	return true
}

func TestLambdaHandler(t *testing.T) {
	tests := []struct {
		name           string
		request        events.APIGatewayProxyRequest
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name: "Valid POST request",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body: `{
					"matrix": {
						"data": [
							[2, 1, 5],
							[1, 1, 3],
							[0, 0, 0]
						]
					}
				}`,
			},
			expectedStatus: 200,
			validateBody: func(t *testing.T, body string) {
				var response RREFResponse
				if err := json.Unmarshal([]byte(body), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				// Check that RREF was calculated
				if response.RREF.Data[0][0] != 1.0 {
					t.Errorf("Expected first element to be 1.0, got %f", response.RREF.Data[0][0])
				}
			},
		},
		{
			name: "OPTIONS request (CORS preflight)",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "OPTIONS",
			},
			expectedStatus: 200,
			validateBody: func(t *testing.T, body string) {
				if body != "" {
					t.Errorf("Expected empty body for OPTIONS, got %s", body)
				}
			},
		},
		{
			name: "Invalid HTTP method",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "GET",
			},
			expectedStatus: 405,
			validateBody: func(t *testing.T, body string) {
				if body == "" {
					t.Error("Expected error message in body")
				}
			},
		},
		{
			name: "Invalid JSON in request body",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       `{"invalid": json}`,
			},
			expectedStatus: 400,
			validateBody: func(t *testing.T, body string) {
				if body == "" {
					t.Error("Expected error message in body")
				}
			},
		},
		{
			name: "Missing matrix in request",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Body:       `{}`,
			},
			expectedStatus: 200, // JSON unmarshal might succeed but matrix will be zero
			validateBody: func(t *testing.T, body string) {
				var response RREFResponse
				if err := json.Unmarshal([]byte(body), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := LambdaHandler(tt.request)
			if err != nil {
				t.Fatalf("LambdaHandler() error = %v", err)
			}

			if response.StatusCode != tt.expectedStatus {
				t.Errorf("LambdaHandler() status = %d, want %d", response.StatusCode, tt.expectedStatus)
			}

			// Check CORS headers
			if response.Headers["Access-Control-Allow-Origin"] != "*" {
				t.Error("Missing CORS header: Access-Control-Allow-Origin")
			}

			if tt.validateBody != nil {
				tt.validateBody(t, response.Body)
			}
		})
	}
}

func TestLambdaHandlerCORSHeaders(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body: `{
			"matrix": {
				"data": [
					[1, 0, 0],
					[0, 1, 0],
					[0, 0, 1]
				]
			}
		}`,
	}

	response, err := LambdaHandler(request)
	if err != nil {
		t.Fatalf("LambdaHandler() error = %v", err)
	}

	// Verify CORS headers
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods":   "",
		"Access-Control-Allow-Headers":  "",
		"Content-Type":                  "application/json",
	}

	if response.Headers["Access-Control-Allow-Origin"] != "*" {
		t.Error("Missing or incorrect Access-Control-Allow-Origin header")
	}

	if response.Headers["Content-Type"] != "application/json" {
		t.Error("Missing or incorrect Content-Type header")
	}

	_ = expectedHeaders // Keep for reference
}

func TestRREFProperties(t *testing.T) {
	// Test that RREF of identity matrix is identity
	identity := [3][3]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	result := RREF(identity)
	if !matricesEqual(result, identity) {
		t.Errorf("RREF of identity should be identity, got %v", result)
	}

	// Test that RREF is idempotent (applying RREF twice gives same result)
	testMatrix := [3][3]float64{
		{2, 1, 5},
		{1, 1, 3},
		{0, 0, 0},
	}
	firstRREF := RREF(testMatrix)
	secondRREF := RREF(firstRREF)
	if !matricesEqual(firstRREF, secondRREF) {
		t.Errorf("RREF should be idempotent, but RREF(RREF(matrix)) != RREF(matrix)")
	}
}

