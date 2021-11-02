# go-crud-apis
- Build Golang Crud Apis using Mux and Postgres.
- Build Golang Crud Apis using GRPC and Postgres.

Endpoints: 

| Method        | Endpoint      |  Description                                  |                   
| ------------- | ------------- | --------------------------------------------- |
| POST          | /signin       | Create a user                                 |
| POST          | /login        | Verify user details and generate jwt token    |
| GET           | /users        | Get list of all users                         |
| GET           | /users/{id}   | Get user by id                                |
| PUT           | /users/{id}   | Update user by id                             |
| DELETE        | /users/{id}   | Delete user by id                             |
