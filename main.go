package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Quiz struct {
	Quiz_id      int       `json:"quiz_id"`
	Quiz_name    string    `json:"quiz_name"`
	Submitted_by string    `json:"submitted_by"`
	Submitted_at time.Time `json:"submitted_at"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)

	}
	defer conn.Close(context.Background())

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		var count string

		err = conn.QueryRow(context.Background(), "select count(*) from \"Quizzes\"").Scan(&count)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)

		}

		fmt.Println("No of questins = ", count)

		c.JSON(200, gin.H{"working": true})
	})

	r.GET("/quiz", func(c *gin.Context) { //Get quiz details and questions based on quiz name

		quiz_name := c.Query("quiz_name")

		var quiz Quiz
		err = conn.QueryRow(
			context.Background(),
			`SELECT * FROM "Quizzes" WHERE "quiz_name"=$1`,
			quiz_name,
		).Scan(&quiz.Quiz_id, &quiz.Quiz_name, &quiz.Submitted_by, &quiz.Submitted_at)
		if err != nil {
			fmt.Println("Query error", err)
			c.JSON(500, gin.H{"error": "Quiz not found or query failed"})
			return
		}

		c.JSON(200, gin.H{"quiz": quiz})

	})

	r.Run(":8000")

}
