package main

import (
    //"log"
   // "net/http"
    "github.com/gin-gonic/gin"

    rdb "github.com/e-dard/netbug"
)

func main() {

    r := gin.Default()
	//r.Use(logging.AccessLogger(os.Stdout))

	//r.Static("/grafana", "../monitor/grafana/")
	//r.Use(static.Serve("/", static.LocalFile("../webexample/refluxjs-todo", false)))
	// This handler will match /user/john but will not match neither /user/ or /user

	//comments
	r.GET("/comment", func(c *gin.Context) {
		c.JSON(200,gin.H{"author":"tsingson","text":"comments"})
	})


    //r := http.NewServeMux()
    rdb.RegisterHandler("/myroute/", r)
    //if err := http.ListenAndServe(":8080", r); err != nil {
    //    log.Fatal(err)
   // }

    r.Run(":8080")
}