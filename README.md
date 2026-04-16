# API Integration and Data Processing

A Go-based API that integrates with the Genderize API to classify names and determine prediction confidence.

## Features
- **Fast Processing:** Built with Go's net/http for <500ms internal response time.
- **Resilient:** Uses a custom HTTP client with a 5-second timeout to prevent upstream hanging.
- **Strict Validation:** Handles missing parameters and non-string inputs with appropriate HTTP codes.

## Setup & Local Development
1. Clone the repo: `git clone <your-repo-link>`
2. Run the server: `go run main.go`
3. The server will be live at `http://localhost:8080`

## API Endpoint
- **GET** `/api/classify?name={name}`