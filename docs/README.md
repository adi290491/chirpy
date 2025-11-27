# Chirpy API Documentation

This document provides a detailed overview of the Chirpy API endpoints.

## Admin

### GET /admin/metrics

- **Description:** Retrieves the number of hits to the file server.
- **Method:** `GET`
- **Path:** `/admin/metrics`
- **Responses:**
  - `200 OK`: Returns an HTML page with the number of hits.

### POST /admin/reset

- **Description:** Resets the file server hits to zero.
- **Method:** `POST`
- **Path:** `/admin/reset`
- **Responses:**
  - `200 OK`: Successfully reset the counter.

## API

### GET /api/healthz

- **Description:** A health check endpoint to verify if the service is running.
- **Method:** `GET`
- **Path:** `/api/healthz`
- **Responses:**
  - `200 OK`: Returns "OK".

### POST /api/users

- **Description:** Creates a new user.
- **Method:** `POST`
- **Path:** `/api/users`
- **Request Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "password"
  }
  ```
- **Responses:**
  - `201 Created`: Returns the newly created user object.
  - `500 Internal Server Error`: If there's an issue creating the user.

### POST /api/login

- **Description:** Logs in a user and returns a JWT.
- **Method:** `POST`
- **Path:** `/api/login`
- **Request Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "password"
  }
  ```
- **Responses:**
  - `200 OK`: Returns a user object with a JWT token.
  - `401 Unauthorized`: If the email or password is incorrect.
  - `500 Internal Server Error`: If there's a server-side issue.

### PUT /api/users

- **Description:** Updates a user's email and password.
- **Method:** `PUT`
- **Path:** `/api/users`
- **Authentication:** Requires a valid JWT in the `Authorization` header.
- **Request Body:**
  ```json
  {
    "email": "new-email@example.com",
    "password": "new-password"
  }
  ```
- **Responses:**
  - `200 OK`: Returns the updated user object.
  - `401 Unauthorized`: If the JWT is missing or invalid.
  - `500 Internal Server Error`: If there's an issue updating the user.

### POST /api/refresh

- **Description:** Refreshes an expired JWT using a refresh token.
- **Method:** `POST`
- **Path:** `/api/refresh`
- **Authentication:** Requires a valid refresh token in the `Authorization` header.
- **Responses:**
  - `200 OK`: Returns a new JWT.
  - `401 Unauthorized`: If the refresh token is invalid or expired.
  - `500 Internal Server Error`: If there's an issue generating a new token.

### POST /api/revoke

- **Description:** Revokes a refresh token.
- **Method:** `POST`
- **Path:** `/api/revoke`
- **Authentication:** Requires a valid refresh token in the `Authorization` header.
- **Responses:**
  - `204 No Content`: The token was successfully revoked.
  - `401 Unauthorized`: If the refresh token is invalid.
  - `500 Internal Server Error`: If there's an issue revoking the token.

### POST /api/chirps

- **Description:** Creates a new chirp.
- **Method:** `POST`
- **Path:** `/api/chirps`
- **Authentication:** Requires a valid JWT in the `Authorization` header.
- **Request Body:**
  ```json
  {
    "body": "This is a new chirp!"
  }
  ```
- **Responses:**
  - `201 Created`: Returns the newly created chirp.
  - `400 Bad Request`: If the chirp is too long.
  - `401 Unauthorized`: If the JWT is missing or invalid.
  - `500 Internal Server Error`: If there's an issue creating the chirp.

### GET /api/chirps

- **Description:** Retrieves all chirps. Can be filtered by `author_id`.
- **Method:** `GET`
- **Path:** `/api/chirps`
- **Query Parameters:**
  - `author_id` (optional): The UUID of the author to filter by.
- **Responses:**
  - `200 OK`: Returns an array of chirps.
  - `500 Internal Server Error`: If there's an issue fetching the chirps.

### GET /api/chirps/{chirpID}

- **Description:** Retrieves a single chirp by its ID.
- **Method:** `GET`
- **Path:** `/api/chirps/{chirpID}`
- **Responses:**
  - `200 OK`: Returns the chirp object.
  - `404 Not Found`: If the chirp with the given ID doesn't exist.
  - `500 Internal Server Error`: If there's an issue parsing the chirp ID.

### DELETE /api/chirps/{chirpID}

- **Description:** Deletes a chirp by its ID.
- **Method:** `DELETE`
- **Path:** `/api/chirps/{chirpID}`
- **Authentication:** Requires a valid JWT in the `Authorization` header. The authenticated user must be the author of the chirp.
- **Responses:**
  - `204 No Content`: The chirp was successfully deleted.
  - `401 Unauthorized`: If the JWT is missing or invalid.
  - `403 Forbidden`: If the user is not the author of the chirp.
  - `404 Not Found`: If the chirp with the given ID doesn't exist.
  - `500 Internal Server Error`: If there's an issue deleting the chirp.

### POST /api/polka/webhooks

- **Description:** A webhook endpoint for Polka to notify of user upgrades.
- **Method:** `POST`
- **Path:** `/api/polka/webhooks`
- **Authentication:** Requires a valid API key in the `Authorization` header.
- **Request Body:**
  ```json
  {
    "event": "user.upgraded",
    "data": {
      "user_id": "user-uuid"
    }
  }
  ```
- **Responses:**
  - `204 No Content`: The webhook was successfully processed.
  - `401 Unauthorized`: If the API key is missing or invalid.
  - `404 Not Found`: If the user is not found.
  - `500 Internal Server Error`: If there's an issue processing the webhook.
