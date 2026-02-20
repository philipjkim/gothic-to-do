package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/philipjkim/gothic-to-do/internal/handler"
	"github.com/philipjkim/gothic-to-do/internal/store"
)

func main() {
	r := gin.Default()

	// Static file serving
	r.Static("/static", "./static")

	// In-memory store & handler
	todoStore := store.NewTodoStore()
	todoHandler := handler.NewTodoHandler(todoStore)

	// Routes
	r.GET("/", todoHandler.Index)
	r.POST("/todos", todoHandler.Create)
	r.PATCH("/todos/:id/toggle", todoHandler.Toggle)
	r.PUT("/todos/:id", todoHandler.Update)
	r.DELETE("/todos/:id", todoHandler.Delete)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
