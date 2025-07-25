# Go API for Quizzo App

An API written in Go using Gin with pgx for PostgreSQL database
Used for Quizzo Web App

# PostgreSQL Schema

<p align="center">
    <img src="quizzo_schema.png">
</p>

# API Reference

API supports 2 functional routes,

```
GET / - Only for testing
```

```
GET /quiz?quiz_name=<QUIZ NAME>
```

```
POST /new_quiz
```

The route `/quiz` uses the query parameter `quiz_name` to get the quiz details and questions from the PostgreSQL DB.
Returns error 404 if quiz is not found and 500 for DB error.

The route `/new_quiz` is to POST a new quiz to the DB.
Returns 400 status code if quiz name is taken.
