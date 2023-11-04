# Shorty
Simple URL Shortener written in golang

## Features:
1. Authentication
2. URL Shortening
3. Protected Short URLs
4. Custom Shortkey
5. Rate Limiting

## Endpoints

- **Endpoint**: `POST /signup`
- **Description**: Register a new user.
- **Request Body**:
  - `email` (string): User's email address.
  - `first_name` (string): User's first name.
  - `last_name` (string): User's last name.
  - `password` (string): User's password.
  - `phone` (string): User's phone number.
- **Response**: The user is registered and receives a token.
  
<hr>

- **Endpoint**: `POST /login`
- **Description**: Authenticate and log in a user.
- **Request Body**:
  - `email` (string): User's email address.
  - `password` (string): User's password.
- **Response**: The user is logged in and receives a token.

<hr>

- **Endpoint**: `POST /shorten`
- **Description**: Shorten a URL. Requires authentication.
- **Request Body**:
  - `url` (string): The original URL to be shortened.
  - `password` (string): An optional password to access the shortened URL.
  - `customShortKey` (string, optional): A custom short key for the shortened URL.
- **Response**: The original URL is shortened, and the shortened URL is provided.

<hr>

- **Endpoint**: `GET /{shortURL}`
- **Description**: Redirect to the original URL associated with the short URL.
- **Query Parameter**:
  - `password` (string, optional): If the original URL requires a password, it should be provided here.
- **Response**: Redirects the user to the original URL if authorized.

<hr>





