# Go Auth Service

Service written in golang which provides a simple authentication flow.

The whole flow looks like:

1. User creates a new account by providing full name, email and password
2. User can sign in by using credentials given in the registration process if credentials are valid then he obtains a pair of tokens, access and refresh tokens. The first one is a short lived token and it is to validate potential user in protected endpoints. The second one is stored in the database and it is used when user want to refresh an access token. All tokens are stored in cookies and those are `httpOnly`.
3. User can make requests to protected endpoints.
4. Client side is responsible for refreshing an access token using refresh token

## How to run

If you would like to run the service locally you must:

- copy `.env.template` file by using `cp .env.template .env` command and fill all required variables.
- Run the `go run cmd/api/main.go` command to run the API application. If you want to have a live reload then you must have `air` package installed. After that, you can simply run `air` command.
