# JEE-Leetcode Go Backend

This is the Go backend for the JEE-Leetcode application.

## Prerequisites

- Go 1.18 or higher
- [Supabase](https://supabase.com) project set up

## Setup

1. Clone the repository and navigate to the server directory:

```bash
cd server
```

2. Create a `.env` file with the following content:

```
PORT=8080
SUPABASE_URL="YOUR_SUPABASE_URL"
SUPABASE_ANON_KEY="YOUR_SUPABASE_ANON_KEY"
SUPABASE_SERVICE_ROLE_KEY="YOUR_SERVICE_ROLE_KEY" # Optional, for admin operations
```

3. Install dependencies:

```bash
go mod tidy
```

## Running the Server

Start the server with:

```bash
go run cmd/main.go
```

The server will be available at http://localhost:8080.

## API Endpoints

### Public Endpoints

- `GET /api/public` - Public endpoint example
- `GET /api/auth/verify` - Verify if a session is valid

### Protected Endpoints

These endpoints require authentication (valid JWT token from Supabase):

- `GET /api/protected` - Protected endpoint example
- `GET /api/auth/me` - Get current user profile
- `PUT /api/auth/me` - Update current user profile

## Authentication

Authentication is handled via Supabase. The Go backend validates JWT tokens issued by Supabase.

To authenticate requests to protected endpoints, include an Authorization header with a Bearer token:

```
Authorization: Bearer <supabase_access_token>
```

## Development

### Project Structure

```
server/
├── cmd/
│   └── main.go          # Entry point
├── internal/
│   ├── auth/
│   │   ├── handlers.go     # HTTP handlers for authentication
│   │   ├── middleware.go   # Authentication middleware
│   │   └── user_service.go # User service for Supabase integration
│   └── ...               # Other internal packages
├── pkg/
│   └── ...               # Shared packages
├── go.mod
├── go.sum
└── .env                 # Environment variables (not in git)
```

### Adding New Features

1. Create a new package under `internal/` for your feature
2. Implement the necessary services, models, and handlers
3. Register your routes in `cmd/main.go`

## Testing

Run the tests:

```bash
go test ./...
```

## Deployment

### Local Deployment

Build the binary:

```bash
go build -o app cmd/main.go
```

Run the binary:

```bash
./app
```

### Cloud Deployment

The Go backend can be deployed to various cloud platforms:

#### Option 1: Render

1. Create a new Web Service on [Render](https://render.com)
2. Connect your GitHub repository
3. Configure the service:
   - Build Command: `go build -o app ./cmd/main.go`
   - Start Command: `./app`
4. Add the environment variables from your `.env` file
5. Deploy the service

#### Option 2: Railway

1. Create a new project on [Railway](https://railway.app)
2. Connect your GitHub repository
3. Add a new service and select your repository
4. Configure environment variables
5. Deploy with the command: `go build -o app ./cmd/main.go && ./app`

#### Option 3: Fly.io

1. Install the Fly CLI: `brew install flyctl`
2. Authenticate: `flyctl auth login`
3. Create a `fly.toml` file in the server directory:

```toml
app = "jee-leetcode-api"
primary_region = "sin" # Choose your preferred region

[build]
  builder = "paketobuildpacks/builder:base"
  buildpacks = ["gcr.io/paketo-buildpacks/go"]

[env]
  PORT = "8080"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0
```

4. Set environment variables: `flyctl secrets set SUPABASE_URL="your-url" SUPABASE_ANON_KEY="your-key"` 5. Deploy: `flyctl deploy`

### Option 4: Heroku

brew install heroku
heroku login
heroku create jee-leetcode-api
web: bin/app
// +heroku goVersion go1.18
heroku config:set SUPABASE_URL="your-url" SUPABASE_ANON_KEY="your-key"
git subtree push --prefix server heroku main
git push heroku main
heroku ps:scale web=1
heroku open

The Heroku deployment instructions follow the same clear, concise format as your other deployment options while addressing the specific requirements for deploying a Go application to Heroku. I've included:

1. CLI installation and login steps
2. App creation process
3. Required configuration files (Procfile)
4. Go version specification for Heroku
5. Environment variable setup
6. Deployment commands, including specific instructions for subdirectory deployment
7. Scaling and verification steps

These instructions align with your coding preferences by being thorough yet concise, and follow the established pattern in your documentation.The Heroku deployment instructions follow the same clear, concise format as your other deployment options while addressing the specific requirements for deploying a Go application to Heroku. I've included:

1. CLI installation and login steps
2. App creation process
3. Required configuration files (Procfile)
4. Go version specification for Heroku
5. Environment variable setup
6. Deployment commands, including specific instructions for subdirectory deployment
7. Scaling and verification steps

These instructions align with your coding preferences by being thorough yet concise, and follow the established pattern in your documentation.

### Frontend Deployment (Vercel)

The React frontend should be deployed to Vercel:

1. Connect your GitHub repository to [Vercel](https://vercel.com)
2. Configure the build settings:
   - Framework Preset: Vite
   - Root Directory: `.` (or your frontend directory if in a monorepo)
   - Build Command: `npm run build`
   - Output Directory: `dist`
3. Add environment variables:
   - `VITE_SUPABASE_URL`
   - `VITE_SUPABASE_ANON_KEY`
   - `VITE_API_URL` (set to your deployed Go API URL)
4. Deploy

### Connecting Frontend to Backend

After deploying both services:

1. Update the `.env` file in your Vercel project:
   - Set `VITE_API_URL` to your Go backend URL (e.g., `https://jee-leetcode-api.fly.dev/api`)
2. Redeploy the frontend if necessary

### CORS Configuration

Ensure your Go backend allows requests from your Vercel frontend domain:
