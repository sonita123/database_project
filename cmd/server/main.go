package main

import (
	"log"
	"net/http"
	"unibazar/project/internal/db"
	"unibazar/project/internal/handlers"
	"unibazar/project/internal/router"
)

func main() {
	//cwd, _ := os.Getwd()
	//fmt.Println("Working directory:", cwd)
	handlers.LoadTemplates("C:/Users/Mahdi Hujat/Desktop/unibazar-go-sql-server-project-main")
	err := db.Connect()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("Database connected")
	r := router.SetupRouter()

	log.Println("Server running on :8080")

	http.ListenAndServe(":8080", r)
}
