# Row Reduction Calculator

A web application that calculates the Reduced Row Echelon Form (RREF) of a 3×3 matrix. Built with React frontend and Go backend.

## Features

- Input a 3×3 matrix through an intuitive web interface
- Calculate RREF using a custom algorithm (no external math libraries)
- Clean, modern UI with real-time input validation

## Project Structure

```
Row-Reduction/
├── backend/
│   ├── main.go          # AWS Lambda handler (build tag: lambda)
│   ├── rowreduce.go     # HTTP server for local development (build tag: !lambda)
│   ├── types.go         # Shared types and RREF algorithm (no build tag)
│   ├── main_test.go     # Test suite for Lambda function
│   ├── build.ps1        # PowerShell build script for Lambda (Windows)
│   └── test.ps1         # PowerShell test script
├── frontend/
│   ├── public/
│   │   └── index.html
│   ├── src/
│   │   ├── App.js       # Main React component
│   │   ├── App.css      # Styling
│   │   ├── index.js     # React entry point
│   │   └── index.css    # Global styles
│   └── package.json
├── go.mod               # Go module file
├── template.yaml        # AWS SAM template for deployment
├── build-lambda.sh      # Build script for Lambda (Linux/Mac)
├── build-lambda.bat     # Build script for Lambda (Windows)
├── test.ps1             # Test script (run from project root)
└── README.md
```

## Setup Instructions

### Local Development

#### Backend (Go HTTP Server)

1. Navigate to the project root directory
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run the Go server:
   ```bash
   # The rowreduce.go file is used for local development (has !lambda build tag)
   go run backend/rowreduce.go backend/types.go
   ```
   The server will start on `http://localhost:8080`

   **Note**: The code uses build tags to separate Lambda and HTTP server code:
   - `main.go` - Lambda handler (only built with `-tags lambda`)
   - `rowreduce.go` - HTTP server (only built without lambda tag)
   - `types.go` - Shared code (always included)

#### Backend (AWS Lambda)

The backend can be deployed to AWS Lambda for serverless execution. See the [AWS Lambda Deployment](#aws-lambda-deployment) section below.

### Frontend (React)

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm start
   ```
   The app will open in your browser at `http://localhost:3000`

## Usage

1. Enter values into the 3×3 matrix grid
2. Click "Calculate RREF" to compute the reduced row echelon form
3. The result will be displayed below the input matrix
4. Use "Clear" to reset the matrix

## API Endpoint

**POST** `/api/rref`

Request body:
```json
{
  "matrix": {
    "data": [
      [1, 2, 3],
      [4, 5, 6],
      [7, 8, 9]
    ]
  }
}
```

Response:
```json
{
  "original": {
    "data": [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
  },
  "rref": {
    "data": [[1, 0, -1], [0, 1, 2], [0, 0, 0]]
  }
}
```

## Algorithm

The RREF algorithm is implemented from scratch in Go without using any external math libraries. It follows the standard Gaussian elimination process:

1. Find the pivot (leading entry) in each row
2. Swap rows if necessary to bring the pivot to the correct position
3. Normalize the pivot row (make the leading entry 1)
4. Eliminate entries above and below the pivot
5. Repeat for all rows

## Testing

The project includes a comprehensive test suite for the Lambda function.

### Running Tests

**Option 1: Using the test script (recommended)**
```powershell
# From project root
.\test.ps1

# Or from backend directory
cd backend
.\test.ps1
```

**Option 2: Manual command**
```powershell
# Make sure GOOS and GOARCH are not set (they interfere with tests)
$env:GOOS = $null; $env:GOARCH = $null; go test -tags lambda -v ./backend
```

**Option 3: From backend directory**
```powershell
cd backend
$env:GOOS = $null; $env:GOARCH = $null; go test -tags lambda -v
```

### Test Coverage

The test suite includes:
- RREF algorithm tests (identity matrix, zero matrix, singular matrices, etc.)
- Lambda handler tests (valid requests, CORS, error handling)
- Mathematical property verification (idempotency, etc.)

**Note**: If you get an error like "not a valid Win32 application", it means `GOOS` and `GOARCH` environment variables are still set from a previous build. Clear them or use the test script.

## AWS Lambda Deployment

The Go backend can be deployed as an AWS Lambda function. The build process uses build tags to separate Lambda and HTTP server code.

### Method 1: Manual Deployment (Windows - Recommended)

1. **Build the Lambda function**:
   ```powershell
   cd backend
   .\build.ps1
   ```
   This will create `rref-api.zip` in the backend directory.

2. **Create Lambda function in AWS Console**:
   - Go to AWS Lambda Console
   - Create a new function
   - Choose "Author from scratch"
   - Runtime: `Custom runtime on Amazon Linux 2023` or `provided.al2023`
   - Architecture: **arm64** (the build script compiles for ARM64)
   - Upload `backend/rref-api.zip`

3. **Configure the function**:
   - Handler: `bootstrap`
   - Timeout: 30 seconds
   - Memory: 128 MB

4. **Set up API Gateway**:
   - Create a new REST API
   - Create a POST method for `/api/rref`
   - Create an OPTIONS method for CORS
   - Integrate both with your Lambda function
   - Deploy the API

5. **Update frontend**: Update the API endpoint in `frontend/src/App.js` to point to your API Gateway URL.

### Method 2: Using AWS SAM

1. **Install AWS SAM CLI** (if not already installed):
   ```bash
   # Follow instructions at: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html
   ```

2. **Build and deploy**:
   ```bash
   # Build the Lambda function
   make build-lambda
   # Or manually:
   cd backend && GOOS=linux GOARCH=arm64 go build -tags lambda -o bootstrap .
   
   # Deploy with SAM
   sam build
   sam deploy --guided
   ```

3. **Update frontend API endpoint**: After deployment, SAM will output the API Gateway URL. Update the frontend to use this URL instead of `localhost:8080`.

### Method 3: Alternative Build Scripts

**On Linux/Mac:**
```bash
./build-lambda.sh
```

**On Windows (alternative):**
```powershell
.\build-lambda.bat
```

**Note**: The `build.ps1` script in the backend directory is the recommended method for Windows as it properly handles environment variables and creates the deployment package automatically.

### Environment Variables

No environment variables are required for basic functionality.

### CORS Configuration

CORS is handled in the Lambda function code. The API Gateway should also be configured to allow CORS if using the REST API.

## Build Tags

The project uses Go build tags to separate Lambda and HTTP server code:

- `//go:build lambda` - Files included only when building for Lambda (`main.go`)
- `//go:build !lambda` - Files included only for local HTTP server (`rowreduce.go`)
- No build tag - Files always included (`types.go`, `main_test.go`)

This allows the same codebase to be used for both local development and Lambda deployment.

## GitHub Pages Deployment

The frontend can be deployed to GitHub Pages while using AWS Lambda for the backend API.

### Quick Setup

1. **Deploy your Lambda function** to AWS (see [AWS Lambda Deployment](#aws-lambda-deployment))
2. **Get your API Gateway URL** from AWS
3. **Enable GitHub Pages** in your repository settings (Settings → Pages → Source: GitHub Actions)
4. **Configure the API URL**:
   - Option A: Add `REACT_APP_API_URL` as a GitHub Secret with your API Gateway URL
   - Option B: Edit `.github/workflows/deploy-pages.yml` and set the API URL directly
5. **Push to main branch** - GitHub Actions will automatically build and deploy

See [GITHUB_PAGES_SETUP.md](GITHUB_PAGES_SETUP.md) for detailed instructions.

### How It Works

- **Frontend**: Static React app deployed to GitHub Pages
- **Backend**: AWS Lambda function (deployed separately)
- **API Communication**: Frontend calls Lambda via API Gateway URL
- **Local Development**: Uses proxy to `localhost:8080` when `REACT_APP_API_URL` is not set

## Technologies

- **Backend**: Go (Golang) - Deployable as AWS Lambda (ARM64) or HTTP server
- **Frontend**: React - Deployable to GitHub Pages or any static hosting
- **No external math libraries** - all calculations are implemented manually
- **Deployment**: 
  - Backend: AWS Lambda with API Gateway
  - Frontend: GitHub Pages (or any static hosting)
- **Testing**: Comprehensive test suite with Go's testing framework
- **CI/CD**: GitHub Actions for automated deployment

