##*-Task Management API-*##

Overview

Task Management API is a simple RESTful API that allows users to manage their tasks. It includes user authentication, task creation, retrieval, updating, and deletion. The API is built using Go, utilizes Redis for rate limiting, and is deployed on Render.com. The project is dockerized and implements CI/CD.

Technologies Used

Programming Language: Go

Framework: Gin

Database: Online PostgreSQL

Cache: Redis (for rate limiting)

Deployment: Render.com

Containerization: Docker

CI/CD: Implemented

API Endpoints

Authentication

POST /login - Logs in a user and returns a JWT token

POST /register - Registers a new user

Task Management (Protected Routes)

Protected routes require authentication using a JWT token.

POST /tasks - Creates a new task

GET /tasks - Retrieves all tasks

GET /tasks/:id - Retrieves a specific task by ID

PUT /tasks/:id - Updates a task by ID

DELETE /tasks/:id - Deletes a task by ID

Deployment

The API is deployed on Render.com with an online PostgreSQL database and Redis for caching. CI/CD is set up to automate deployments.

Running Locally

Prerequisites

Go installed

Docker installed (optional)

Steps

Clone the repository:

git clone https://github.com/ratheeshkumar25/taskmanagement.git
cd ratheeshkumar25-taskmanagement

Run the API:

go run cmd/main.go

Access the API at http://localhost:8080

Running with Docker

docker build -t task-management-api .
docker run -p 8080:8080 task-management-api

Rate Limiting

The API implements rate limiting using Redis. Each user is allowed 60 requests per minute.

CI/CD

The project has a CI/CD pipeline configured to automate testing and deployment.

License

MIT License