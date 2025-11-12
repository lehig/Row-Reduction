@echo off
REM Build script for AWS Lambda deployment (Windows)

echo Building Lambda function...

REM Set build environment for Linux (Lambda runs on Linux)
set GOOS=linux
set GOARCH=amd64

REM Build the Lambda function
cd backend
go build -tags lambda -o bootstrap lambda.go

REM Create deployment package (requires 7-Zip or similar)
echo.
echo Build complete! File: backend\bootstrap
echo.
echo To create deployment package:
echo 1. Zip the bootstrap file as lambda-deployment.zip
echo 2. Upload to AWS Lambda
echo 3. Set handler to: bootstrap
echo 4. Set runtime to: provided.al2023 (or provided.al2)
echo 5. Configure API Gateway trigger

cd ..

