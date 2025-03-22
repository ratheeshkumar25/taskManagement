Task Management API

Overview

Task Management API is a RESTful service that allows users to manage their tasks. It includes user authentication, task creation, retrieval, updating, and deletion. The API is built using Go, utilizes Redis for rate limiting, and is deployed on Render.com. The project is containerized using Docker and implements CI/CD for automated deployments.

Technologies Used

Programming Language: Go

Framework: Gin

Database: PostgreSQL (Online)

Cache: Redis (for rate limiting)

Deployment: Render.com

Containerization: Docker

CI/CD: Implemented

API Endpoints and Versioning (Base URL: https://task-management-v1-0.onrender.com/api/v1)

Authentication

**POST **/register - Registers a new user

**POST **/login - Logs in a user and returns a JWT token

Task Management (Protected Routes - Requires JWT Authentication)

**POST **/tasks - Creates a new task

**GET **/tasks - Retrieves all tasks

**GET **/tasks/:id - Retrieves a specific task by ID

**PUT **/tasks/:id - Updates a task by ID

**DELETE **/tasks/:id - Deletes a task by ID

Deployment

The API is deployed on Render.com with an online PostgreSQL database and Redis for caching. CI/CD is set up to automate deployments.

Running Locally

Prerequisites

Go installed

Docker installed (optional)

Steps

Clone the repository:

git clone https://github.com/ratheeshkumar25/taskmanagement.git
cd taskmanagement

Run the API:

go run cmd/main.go

Access the API at http://localhost:3000

Running with Docker

Build the Docker image:

docker build -t task-management-api .

Run the container:

docker run -p 3000:3000 task-management-api

Rate Limiting

The API implements rate limiting using Redis. Each user is allowed 60 requests per minute.

Sample API Requests & Responses

Register User

Request

{
    "username": "rahul@gmail.com",
    "password": "revathy123"
}

Response

{
    "error": "failed to register user - username already exists"
}

Login User

Request

{
    "username": "rahul@gmail.com",
    "password": "revathy123"
}

Response

Invalid Credentials:

{
    "error": "invalid credentials"
}

Successful Login:

{
    "token": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}

Create Task

Request

{
    "Title": "Fixing redis config file",
    "Description": "Finish the development working on the deployment.",
    "Status": "Pending",
    "DueDate": "2025-02-25T00:00:00Z"
}

Get All Tasks

Response

{
    "limit": 3,
    "page": 1,
    "tasks": [
        {
            "ID": 2,
            "Title": "Fixing Backend API issue",
            "Status": "Pending",
            "CreatedAt": "2025-03-21T13:11:58Z",
            "UpdatedAt": "2025-03-21T13:11:58Z"
        },
        {
            "ID": 5,
            "Title": "Fixing redis config file",
            "Status": "Pending",
            "CreatedAt": "2025-03-21T18:42:05Z",
            "UpdatedAt": "2025-03-21T18:42:05Z"
        },
        {
            "ID": 6,
            "Title": "Fixing redis config file",
            "Status": "Pending",
            "CreatedAt": "2025-03-21T18:43:01Z",
            "UpdatedAt": "2025-03-21T18:43:01Z"
        }
    ],
    "total": 3
}

Get Task by ID

Response

{
    "id": 2,
    "title": "Fixing Backend API issue",
    "description": "Finish the development of the backend and Preparing Docker file.",
    "status": "Pending",
    "dueDate": "2025-02-25T05:30:00+05:30",
    "createdAt": "2025-03-21T18:41:58.344666+05:30",
    "updatedAt": "2025-03-21T18:41:58.344667+05:30"
}

Update Task

Request

{
    "title": "Update Backend API Logic",
    "description": "Refactor the backend API and fix endpoint bugs",
    "status": "In Progress",
    "dueDate": "2025-03-30T00:00:00Z"
}

Response

{
    "id": 1,
    "title": "Update Backend API Logic",
    "description": "Refactor the backend API and fix endpoint bugs",
    "status": "In Progress",
    "dueDate": "2025-03-30T00:00:00Z",
    "createdAt": "0001-01-01T00:00:00Z",
    "updatedAt": "2025-03-21T19:27:23.508218153+05:30"
}

Delete Task

Response

{
    "message": "task deleted successfully"
}

CI/CD

The project has a CI/CD pipeline configured to automate testing and deployment.

License

MIT License

