
# Simple Dating App

A Simple Dating App project for test and learn purposes only.


## Authors

- [@sigit14ap](https://www.github.com/sigit14ap)


## Documentation

- CMD Folder
The cmd folder typically contains the main executable files for your application.
- Repository folder
Focus on write and read to the table in database.
- Service folder
Functionality for business logic purposes.
## API Reference

#### Create account

```http
  POST /api/signup
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `username` | `string` | **Required**|
| `password` | `string` | **Required**|

#### Login account

```http
  POST /api/login
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `username` | `string` | **Required**|
| `password` | `string` | **Required**|


#### Get next user to swipe

```http
  GET /api/next-user
```

| Parameter |
| :-------- |
| **No parameter required**|


#### Swipe user

```http
  POST /api/swipe
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `target_user_id` | int | **Required**|
| `match` | `bool` | **Required**|

#### Buy unlimited swipe package

```http
  POST /api/buy-package-unlimited-swipe
```
| Parameter |
| :-------- |
| **No parameter required**|



## Setup
- Copy .env.example file
```bash
  cp .env.example .env
```
- Configure your .env file

## Deployment

### Without Docker

```bash
  go mod download
  go run cmd/main.go
```

### With Docker
```bash
docker-compose build
docker-compose up -d
```
- Note :
Please change these variable in docker-compose.yml to same as your .env file
```
MYSQL_ROOT_PASSWORD: password
MYSQL_DATABASE: simple_dating_app
```


## Running Tests

Use Postman to test the APIs

- Import ```[Test] Simple Dating App.postman_collection.json``` file to the Postman.
- Go to variable collection.
- Fill the ```base_url``` current value.

