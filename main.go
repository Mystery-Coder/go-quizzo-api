package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Quiz struct {
	Quiz_id      int       `json:"quiz_id"`
	Quiz_name    string    `json:"quiz_name"`
	Submitted_by string    `json:"submitted_by"`
	Submitted_at time.Time `json:"submitted_at"`
}

type Question struct {
	Question_id int    `json:"question_id"`
	Quiz_id     int    `json:"quiz_id"`
	Question    string `json:"question"`
	Answer      string `json:"answer"`
	Option1     string `json:"option1"`
	Option2     string `json:"option2"`
	Option3     string `json:"option3"`
	Option4     string `json:"option4"`
}

type QuizPost struct {
	Quiz      Quiz       `json:"quiz"`
	Questions []Question `json:"questions"`
}

func main() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL_SUPABASE"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)

	}
	defer conn.Close(context.Background())

	r := gin.Default()
	r.Use(cors.New(cors.Config{AllowOrigins: []string{"https://quizzo-angular.vercel.app", "http://localhost:4200"}, AllowMethods: []string{"GET", "POST"}}))

	r.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{"working": true})
	})

	r.GET("/quiz", func(c *gin.Context) { //Get quiz details and questions based on quiz name

		quiz_name := c.Query("quiz_name")

		var quiz Quiz
		err = conn.QueryRow(
			context.Background(),
			`SELECT * FROM "Quizzes" WHERE "quiz_name"=$1;`,
			quiz_name,
		).Scan(&quiz.Quiz_id, &quiz.Quiz_name, &quiz.Submitted_by, &quiz.Submitted_at)

		if err != nil {
			if err == pgx.ErrNoRows {

				c.JSON(404, gin.H{"error": "Quiz not found"})
			} else {
				fmt.Println("Query error:", err)
				c.JSON(500, gin.H{"error": "Query failed"})
			}
			return
		}

		var questions []Question
		rows, err := conn.Query(context.Background(), `SELECT * from "Questions" WHERE "quiz_id"=$1;`, quiz.Quiz_id)

		if err != nil {
			fmt.Println("Query err", err)
			c.JSON(500, gin.H{"error": "query error"})
			return
		}

		for rows.Next() {
			var q Question

			err := rows.Scan(&q.Question_id, &q.Quiz_id, &q.Question, &q.Answer, &q.Option1, &q.Option2, &q.Option3, &q.Option4)

			if err != nil {
				c.JSON(500, gin.H{"error": "query error"})
				return
			}

			questions = append(questions, q)
		}

		c.JSON(200, gin.H{"quiz": quiz, "questions": questions})

	})

	r.GET("/quiz_exists", func(c *gin.Context) {
		quiz_name := c.Query("quiz_name")

		var exists bool
		err := conn.QueryRow(context.Background(), `SELECT EXISTS (SELECT 1 FROM "Quizzes" where "quiz_name"=$1)`, quiz_name).Scan(&exists)

		if err != nil {
			c.JSON(500, gin.H{"error": "Query Error"})
			return
		}
		c.JSON(200, gin.H{"exists": exists})
	})

	r.POST("/quiz", func(c *gin.Context) {
		/*
			INSERT INTO "Quizzes" ("quiz_name", "submitted_by") VALUES ('PhyTestQ', 'Test') returning quiz_id;

			INSERT INTO "Questions" ("quiz_id", "question", "answer", "option1", "option2", "option3", "option4")
			VALUES (5, 'Symbol of Force?', 'N', 'Pa', 'N', 'J', 's');

			Insert queries for reference

			PostgreSQL auto calculates IDs and time
			Expects data as,
			{
				"quiz" : {
					"quiz_id" : 0 //DUMMY VALUES to reuse Structs
					"quiz_name" : Quiz1,
					"submitted_by" : Name,
					"submitted_at" : 0 //DUMMY VALUES to reuse Structs
				},
				"questions" : [
					{
						"question_id" : 0, //DUMMY VALUES to reuse Structs
						"quiz_id" : 0, //DUMMY VALUES to reuse Structs
						"question" : ".....",
						"answer" : ".....",
						"option1" : ".....",
						"option2" : ".....",
						"option3" : ".....",
						"option4" : ".....",

					},
					{
						...
					}
				]
			}
		*/

		var quizData QuizPost

		if err := c.BindJSON(&quizData); err != nil {
			fmt.Println("Binding error", err)
		}

		// prettyJSON, err := json.MarshalIndent(quizData, "", "  ") // Using two spaces for indentation
		// if err != nil {
		// 	log.Fatalf("Error marshalling JSON: %v", err)
		// }

		// fmt.Println(string(prettyJSON))

		var quiz_id int
		err := conn.QueryRow(context.Background(),
			`INSERT INTO "Quizzes" ("quiz_name", "submitted_by") VALUES ($1, $2) returning "quiz_id";`,
			quizData.Quiz.Quiz_name,
			quizData.Quiz.Submitted_by).Scan(&quiz_id)

		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
				// Unique constraint violation on quiz_name
				c.JSON(400, gin.H{"error": "QuizNameExists"})
			} else {
				// Some other DB error
				c.JSON(500, gin.H{"error": "InsertFailed"})
			}
			return
		}

		questions := quizData.Questions

		for qIdx := 0; qIdx < len(questions); qIdx++ {
			q := questions[qIdx]

			_, err := conn.Exec(context.Background(), `INSERT INTO "Questions" ("quiz_id", "question", "answer", "option1", "option2", "option3", "option4")
			VALUES ($1, $2, $3, $4, $5, $6, $7);`, quiz_id, q.Question, q.Answer, q.Option1, q.Option2, q.Option3, q.Option4)

			if err != nil {
				c.JSON(500, gin.H{"error": "InsertError"})
				return
			}
		}

		c.JSON(200, gin.H{"success": true})

	})

	r.Run(":8080")

}
