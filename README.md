# trail-data-service
A Go REST API that manages trail/waypoint data

# Features
# CRUD operations for trails (name, coordinates, difficulty, length)
# Endpoint to find trails within a radius of given coordinates
# Basic authentication middleware
# Input validation and error handling
# Unit tests for key functions
# Docker containerization

# Tech Stack
* Language: Go 1.24+
* Router: chi
* UUIDs: github.com/google/uuid
* Testing: Go standard library

# Project Structure
```
trail-api/
├── main.go
├── handlers/
│   └── trails.go
├── models/
│   └── trail.go
├── middleware/
│   └── auth.go
├── storage/
│   └── memory.go
└── Dockerfile
```

# Running the Service
First run the included build.sh script to build the container images
```
./build.sh
```

Then start the application in docker with the following command
```
docker run -d -p 8080:8080 --name trail-service trail-data-service
```


# Example Usage (cURL)
## POST /trails - create trail
```
curl -X POST http://localhost:8080/trails \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Lamar River Trail",
    "lat": 44.8472,
    "lon": -109.6278,
    "difficulty": "hard",
    "length_km": 53
}'
```

## List all trails
GET /trails/{uid} - list trail by id
```curl http://localhost:8080/trails/6f03765b-6a3d-44df-9c1f-f3341f089c23```

GET /trails - list all trails
```curl http://localhost:8080/trails```

GET /trails/nearby?lat=X&lon=Y&radius-km=Z - proximity search
```curl http:///trails/nearby?lat=44.8472&lon=-109.6278&radius-km=50```

# Design Considerations
* Dependency Injection is used for loose coupling between components.
* Interface-Driven Architecture enables testability and future extensibility (e.g., database-backed repo).
* Validation is handled at the request model level to separate concerns cleanly.
* The service layer enforces any domain-specific business rules.

# Tests
`go test ./...`
Tests cover handler logic, service behavior, and in-memory repo operations.

# Future Improvements / Next Steps
TBD

# Time Spent
TBD

# Author
David Nakolan - david.nakolan@gmail.com
