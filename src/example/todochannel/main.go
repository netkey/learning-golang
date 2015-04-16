package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	//"log"
	. "example/todochannel/todos"
	"fmt"
	"github.com/mgutz/logxi/v1"
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
			//log.Fatal(err)
			fmt.Println("error:", err)
		}
	}

	// create channel to communicate over
	jobs := make(chan Job)

	// start watching jobs channel for work
	go ProcessJobs(jobs, Db)

	// create dependencies
	client := &TodoClient{Jobs: jobs}
	handlers := &TodoHandlers{Client: client}

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
