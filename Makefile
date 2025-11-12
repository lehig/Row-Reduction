.PHONY: build-lambda deploy-local clean

# Build Lambda function for deployment
build-lambda:
	@echo "Building Lambda function..."
	@cd backend && GOOS=linux GOARCH=amd64 go build -tags lambda -o bootstrap lambda.go
	@cd backend && zip -r ../lambda-deployment.zip bootstrap
	@cd backend && rm bootstrap
	@echo "Build complete! Deployment package: lambda-deployment.zip"

# Build for local testing (regular HTTP server)
build-local:
	@echo "Building local server..."
	@go build -o bin/server backend/main.go
	@echo "Build complete! Run with: ./bin/server"

# Deploy using SAM (requires AWS SAM CLI)
deploy:
	@echo "Deploying with SAM..."
	@sam build
	@sam deploy --guided

# Clean build artifacts
clean:
	@rm -f lambda-deployment.zip
	@rm -f backend/bootstrap
	@rm -f bin/server
	@rm -rf .aws-sam

