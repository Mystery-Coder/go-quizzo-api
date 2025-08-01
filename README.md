# Go API for Quizzo App

An API written in Go using Gin with pgx for PostgreSQL database
Used for Quizzo Web App

# PostgreSQL Schema

<p align="center">
    <img src="quizzo_schema.png">
</p>

# API Reference

API supports 3 functional routes,

```http
GET /quiz?quiz_name=<QUIZ NAME>
```

```http
GET /quiz_exists?quiz_name=<QUIZ Name>
```

```http
POST /new_quiz
```

The route `/quiz` uses the query parameter `quiz_name` to get the quiz details and questions from the PostgreSQL DB.
Returns error 404 if quiz is not found and 500 for DB error.

The route `/new_quiz` is to POST a new quiz to the DB.
Returns 400 status code if quiz name is taken.

The `/quiz_exists` route is just an SQL exists statement for uniqueness of a quiz name.
