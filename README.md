# otaQku Tasks
Simple tasks management

## How to run
All you need to do is to create an `.env` file (see `.env.example`) and run these commands:
```sh
make run-stack
make run
```

`make run-stack` will run a local PostgreSQL Docker container, while `make run` will run the API service

## API Documentation
The API documentation is in the form of Postman collection. You can import `Auth.postman_collection.json` for the account registration and authentication API, and import `Task.postman_collection.json` for the task management API