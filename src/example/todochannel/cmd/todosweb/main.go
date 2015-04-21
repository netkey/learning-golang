package main

import (

	// web
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	//"github.com/olebedev/gin-cache"

	// log
	log "github.com/inconshreveable/log15"

	// go standard lib
	"fmt"
	"io/ioutil"
	//"time"

	// project module
	"example/todochannel/todos"

	// debug
	//"runtime/pprof"
	//"flag"
	//"os"
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

	// start web service
	// configure routes
	router := gin.Default()
	// cache ??
	// action
	router.POST("/todo", handlers.AddTodo)
	router.GET("/todo", handlers.GetTodos)
	router.GET("/todo/:id", handlers.GetTodo)
	router.PUT("/todo/:id", handlers.SaveTodo)
	router.DELETE("/todo/:id", handlers.DeleteTodo)

	// debug
	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	ginpprof.Wrapper(router)

	// start web server
	router.Run(":8080")
}
