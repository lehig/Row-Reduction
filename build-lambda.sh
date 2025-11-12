#!/bin/bash

# Build script for AWS Lambda deployment

echo "Building Lambda function..."

# Set build environment for Linux (Lambda runs on Linux)
export GOOS=linux
export GOARCH=amd64

# Build the Lambda function
cd backend
go build -tags lambda -o bootstrap lambda.go

# Create deployment package
zip -r ../lambda-deployment.zip bootstrap

# Clean up
rm bootstrap

echo "Build complete! Deployment package: lambda-deployment.zip"
echo ""
echo "To deploy:"
echo "1. Upload lambda-deployment.zip to AWS Lambda"
echo "2. Set handler to: bootstrap"
echo "3. Set runtime to: provided.al2023 (or provided.al2)"
echo "4. Configure API Gateway trigger"

