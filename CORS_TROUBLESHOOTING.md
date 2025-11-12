# CORS Troubleshooting Guide

If you're getting CORS errors when calling your Lambda API from GitHub Pages, follow these steps:

## Issue 1: Missing `/api/rref` in API URL

**Error**: `403 Forbidden` or `CORS header missing`

**Solution**: Make sure your API Gateway URL includes the full path:
- ❌ Wrong: `https://lfyth34k32.execute-api.us-east-1.amazonaws.com/dev`
- ✅ Correct: `https://lfyth34k32.execute-api.us-east-1.amazonaws.com/dev/api/rref`

Update your GitHub Secret `REACT_APP_API_URL` with the full path.

## Issue 2: API Gateway CORS Configuration

Even though the Lambda function returns CORS headers, API Gateway also needs to be configured.

### For REST API (API Gateway v1):

1. Go to API Gateway Console → Your API → Resources
2. Select `/api/rref` resource
3. Click "Actions" → "Enable CORS"
4. Configure:
   - **Access-Control-Allow-Origin**: `*` (or your specific domain)
   - **Access-Control-Allow-Headers**: `Content-Type`
   - **Access-Control-Allow-Methods**: `POST, OPTIONS`
5. Click "Enable CORS and replace existing CORS headers"
6. **Important**: Deploy your API (Actions → Deploy API)

### For HTTP API (API Gateway v2):

1. Go to API Gateway Console → Your API → Routes
2. Select your route (e.g., `POST /api/rref`)
3. Go to "CORS" tab
4. Configure:
   - **Access-Control-Allow-Origin**: `*`
   - **Access-Control-Allow-Headers**: `Content-Type`
   - **Access-Control-Allow-Methods**: `POST, OPTIONS`
5. Save and redeploy

### For OPTIONS Method:

Make sure you have an OPTIONS method configured:
1. In API Gateway, add an OPTIONS method to `/api/rref`
2. Set it to return 200 with CORS headers
3. Or configure it to use the Lambda function (which handles OPTIONS)

## Issue 3: Testing Your API

Test your API directly to verify it works:

```bash
# Test POST request
curl -X POST https://lfyth34k32.execute-api.us-east-1.amazonaws.com/dev/api/rref \
  -H "Content-Type: application/json" \
  -d '{"matrix":{"data":[[1,2,3],[4,5,6],[7,8,9]]}}'

# Test OPTIONS (CORS preflight)
curl -X OPTIONS https://lfyth34k32.execute-api.us-east-1.amazonaws.com/dev/api/rref \
  -H "Origin: https://lehig.github.io" \
  -H "Access-Control-Request-Method: POST" \
  -v
```

The OPTIONS request should return CORS headers.

## Quick Fix Checklist

- [ ] API URL includes full path: `/dev/api/rref` (not just `/dev`)
- [ ] API Gateway CORS is enabled for the resource
- [ ] OPTIONS method is configured (or handled by Lambda)
- [ ] API is deployed after CORS changes
- [ ] GitHub Secret `REACT_APP_API_URL` has the correct full URL
- [ ] Frontend is rebuilt and redeployed after URL change

## Common Mistakes

1. **Forgetting to deploy API after CORS changes** - CORS changes don't take effect until you deploy
2. **Missing `/api/rref` in URL** - Only the base URL won't work
3. **Not configuring OPTIONS method** - Browsers send OPTIONS preflight requests
4. **Using wrong API Gateway type** - REST API (v1) vs HTTP API (v2) have different CORS settings

