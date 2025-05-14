# Movie Actors API Documentation

## Overview
RESTful API for managing movies, actors, and their relationships with user authentication.

## Base URL
`http://localhost:8080/api`

## Initialization

1. Clone repository:

   ```bash
   git clone git@github.com:enagtim/movielibrary-api.git
   ```

2. Open directory:

   ```bash
   cd movielibrary-api
   ```

3. Create file `.env` in the root of the project:

   ```bash
   touch .env
   ```

4. Open `.env` and add:

   ```
   DB_USER=postgres
   DB_PASSWORD=secret
   DB_NAME=filmlibrary
   JWT_SECRET="/2+XnmJGz1j3ehIVI/5P9kl+CghrE3DcS7rnT+qar5w="
   ```

## Authentication
Requires JWT token for protected endpoints (except `/auth/register` and `/auth/login`).

## Endpoints

### Authentication
| Method | Endpoint       | Description          |
|--------|----------------|----------------------|
| POST   | `/auth/register` | Register new user    |
| POST   | `/auth/login`    | Login and get JWT    |

### Actors
| Method | Endpoint             | Description                         | Role Required |
|--------|----------------------|-------------------------------------|---------------|
| POST   | `/admin/actors`      | Create new actor                    | admin         |
| GET    | `/actors`           | Get all actors with their movies     | user          |
| GET    | `/actors/{id}`      | Get actor by ID                      | user          |
| PUT    | `/admin/actors/{id}`| Fully update actor                   | admin         |
| PATCH  | `/admin/actors/{id}`| Partially update actor               | admin         |
| DELETE | `/admin/actors/{id}`| Delete actor                         | admin         |

### Movies
| Method | Endpoint               | Description                         | Role Required |
|--------|------------------------|-------------------------------------|---------------|
| POST   | `/admin/movies`        | Create new movie                    | admin         |
| GET    | `/movies`             | Get all movies (sortable)            | user          |
| GET    | `/movies/{id}`        | Get movie by ID                      | user          |
| PUT    | `/admin/movies/{id}`  | Fully update movie                   | admin         |
| PATCH  | `/admin/movies/{id}`  | Partially update movie               | admin         |
| DELETE | `/admin/movies/{id}`  | Delete movie                         | admin         |
| GET    | `/movies/search/title`| Search movies by title               | user          |
| GET    | `/movies/search/actorname`| Search movies by actor name      | user          |

## Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1", "password":"pass123", "role":"user"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1", "password":"pass123"}'
```

### Create Actor (Admin)
```bash
curl -X POST http://localhost:8080/api/admin/actors \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name":"Leonardo DiCaprio", "gender":"male", "birth_date":"1974-11-11"}'
```

### Get Movies Sorted
```bash
curl -X GET "http://localhost:8080/api/movies?sortBy=title" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Search Movies
```bash
curl -X GET "http://localhost:8080/api/movies/search/title?title=Inception" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Error Responses
- `400 Bad Request` - Invalid input data
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Models

### User
```json
{
  "username": "string",
  "password": "string",
  "role": "user|admin"
}
```

### Actor
```json
{
  "id": "number",
  "name": "string",
  "gender": "string",
  "birth_date": "string"
}
```

### Movie
```json
{
  "id": "number",
  "title": "string",
  "description": "string",
  "release_date": "string",
  "rating": "number"
}
```