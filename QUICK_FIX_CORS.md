# Quick Fix for CORS Error

## The Problem

Your API URL is missing the `/api/rref` path. The error shows:
- ❌ Current: `https://lfyth34k32.execute-api.us-east-1.amazonaws.com/dev`
- ✅ Should be: `https://lfyth34k32.execute-api.us-east-1.amazonaws.com/dev/api/rref`

## Step-by-Step Fix

### Step 1: Update GitHub Secret (CRITICAL)

1. Go to your GitHub repository
2. Settings → Secrets and variables → Actions
3. Find `REACT_APP_API_URL` secret
4. **Edit it** and make sure it's exactly:
   ```
   https://lfyth34k32.execute-api.us-east-1.amazonaws.com/dev/api/rref
   ```
   ⚠️ **Must include `/api/rref` at the end!**

### Step 2: Trigger a New Build

The environment variable is baked into the build at build time. You MUST rebuild:

**Option A: Push a small change**
```bash
# Make a small change (like adding a comment)
# Then commit and push
git commit --allow-empty -m "Trigger rebuild with updated API URL"
git push origin main
```

**Option B: Manually trigger workflow**
1. Go to your repository → Actions tab
2. Click "Deploy to GitHub Pages" workflow
3. Click "Run workflow" → "Run workflow"
4. Wait for it to complete (2-3 minutes)

### Step 3: Verify API Gateway CORS

1. Go to AWS Console → API Gateway
2. Select your API
3. Click on `/api/rref` resource (or create it if it doesn't exist)
4. Click "Actions" → "Enable CORS"
5. Set:
   - **Access-Control-Allow-Origin**: `*`
   - **Access-Control-Allow-Headers**: `Content-Type`
   - **Access-Control-Allow-Methods**: `POST, OPTIONS`
6. Click "Enable CORS and replace existing CORS headers"
7. **IMPORTANT**: Click "Actions" → "Deploy API" → Select `dev` stage → "Deploy"

### Step 4: Verify OPTIONS Method

1. In API Gateway, check if `/api/rref` has an OPTIONS method
2. If not:
   - Click on `/api/rref`
   - "Actions" → "Create Method" → Select "OPTIONS"
   - Integration type: Mock or Lambda
   - Method Response: Add CORS headers
   - Deploy API

### Step 5: Test

After the new build completes:
1. Clear your browser cache (Ctrl+Shift+Delete)
2. Visit your GitHub Pages site
3. Try calculating RREF again

## Verification Checklist

- [ ] GitHub Secret `REACT_APP_API_URL` = `https://lfyth34k32.execute-api.us-east-1.amazonaws.com/dev/api/rref` (with `/api/rref`)
- [ ] New build triggered and completed
- [ ] API Gateway CORS enabled on `/api/rref`
- [ ] OPTIONS method exists on `/api/rref`
- [ ] API deployed after CORS changes
- [ ] Browser cache cleared

## Still Not Working?

1. **Check the built JavaScript**: 
   - View page source on GitHub Pages
   - Search for your API URL
   - Verify it includes `/api/rref`

2. **Test API directly**:
   ```bash
   curl -X POST https://lfyth34k32.execute-api.us-east-1.amazonaws.com/dev/api/rref \
     -H "Content-Type: application/json" \
     -d '{"matrix":{"data":[[1,2,3],[4,5,6],[7,8,9]]}}'
   ```

3. **Check browser console**: Look for the exact URL being called

4. **Verify API Gateway stage**: Make sure you're using the correct stage (`dev` vs `Prod`)

