package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	//"log"
	"example/todochannel/todos"
	"fmt"
	log "github.com/inconshreveable/log15"
)

const Db = "./todos.json"

var (
	logger log.Logger
)

func main() {

	// initialize empty-object json file if not found
	if _, err := ioutil.ReadFile(Db); err != nil {
		str := "{}"
		if err = ioutil.WriteFile(Db, []byte(str), 0644); err != nil {
			log.Crit("Unable write to file", "error", err)
			fmt.Println("error:", err)
		}
	}

	// create channel to communicate over
	jobs := make(chan todos.Job)

	log.Crit("start process job")
	// start watching jobs channel for work
	go todos.ProcessJobs(jobs, Db)

	// create dependencies
	client := &todos.TodoClient{Jobs: jobs}
	handlers := &todos.TodoHandlers{Client: client}

	// configure routes
	r := gin.Default()

	r.POST("/todo", handlers.AddTodo)
	r.GET("/todo", handlers.GetTodos)
	r.GET("/todo/:id", handlers.GetTodo)
	r.PUT("/todo/:id", handlers.SaveTodo)
	r.DELETE("/todo/:id", handlers.DeleteTodo)

	// start web server
	r.Run(":8080")
}
