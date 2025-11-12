# GitHub Pages Setup Guide

This guide explains how to deploy the frontend to GitHub Pages while using AWS Lambda for the backend.

## Prerequisites

1. Your Lambda function deployed to AWS with API Gateway
2. Your API Gateway URL (e.g., `https://abc123.execute-api.us-east-1.amazonaws.com/Prod/api/rref`)
3. GitHub repository with GitHub Pages enabled

## Setup Steps

### 1. Enable GitHub Pages

1. Go to your repository Settings
2. Navigate to Pages
3. Under "Source", select "GitHub Actions"

### 2. Configure API URL

You have two options:

**Option A: Use GitHub Secrets (Recommended)**
1. Go to repository Settings → Secrets and variables → Actions
2. Add a new secret named `REACT_APP_API_URL`
3. Set the value to your API Gateway URL (e.g., `https://abc123.execute-api.us-east-1.amazonaws.com/Prod/api/rref`)

**Option B: Update the workflow file directly**
1. Edit `.github/workflows/deploy-pages.yml`
2. Replace the placeholder URL in the `REACT_APP_API_URL` environment variable with your actual API Gateway URL

### 3. Deploy

1. Push your code to the `main` branch
2. The GitHub Actions workflow will automatically:
   - Build the React app with the API URL configured
   - Deploy it to GitHub Pages

### 4. Access Your Site

Your site will be available at:
- `https://<your-username>.github.io/<repository-name>/`

## Local Development

For local development, the app will use the proxy configured in `package.json` (localhost:8080) when `REACT_APP_API_URL` is not set.

To test with your Lambda API:
1. Create a `.env` file in the `frontend` directory
2. Add: `REACT_APP_API_URL=https://your-api-gateway-url.execute-api.region.amazonaws.com/Prod/api/rref`
3. Restart the development server

## Troubleshooting

### CORS Issues

If you see CORS errors, check:

1. **API URL includes full path**: 
   - ❌ Wrong: `https://abc123.execute-api.us-east-1.amazonaws.com/dev`
   - ✅ Correct: `https://abc123.execute-api.us-east-1.amazonaws.com/dev/api/rref`

2. **API Gateway CORS Configuration**:
   - Go to API Gateway Console
   - Enable CORS on your `/api/rref` resource
   - Make sure OPTIONS method is configured
   - **Deploy your API** after making CORS changes

3. **See [CORS_TROUBLESHOOTING.md](CORS_TROUBLESHOOTING.md) for detailed instructions**

### API Not Working

1. Verify your API Gateway URL includes the full path: `/dev/api/rref` or `/Prod/api/rref`
2. Test the API directly using curl or Postman
3. Check browser console for error messages
4. Verify the GitHub Secret `REACT_APP_API_URL` has the complete URL

## Notes

- The frontend is a static React app, perfect for GitHub Pages
- The backend must be deployed separately (AWS Lambda)
- GitHub Pages only hosts static files - no server-side code runs on GitHub

