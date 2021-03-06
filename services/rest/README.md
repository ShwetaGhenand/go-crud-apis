# REST APIs
Build Golang Crud Apis using Mux and Postgres.
User Service is REST API server implementation using gorila/mux.

## Run locally

- Clone the repository
```
git clone git@github.com:ShwetaGhenand/go-crud-apis.git
```
-  Open a terminal in the root directory of your code and run the following command to start the application.

### Run with local go installtion
```
go run main server
```

### Run with Docker
1. Start the container in the background
```
make up
```
2. See logs
```
make logs
```
3. Stop container 
```
make down
```

### Testing the API

Test REST API using 
```
curl --location --request GET 'http://localhost:8081/health'
```
