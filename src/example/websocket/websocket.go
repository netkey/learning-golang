package main

import (
	"fmt"
	"net/http"
	"os"
	//"runtime/debug"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var dir string
var port int
var indexs []string

func init() {
	indexs = []string{"index.html", "index.htm"}
}

func main() {

	r := gin.Default()
	r.LoadHTMLFiles("../public/websocket/index.html")

	r.GET("/websocket", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// websocket
	r.GET("/ws", func(c *gin.Context) {
		webSocketHandler(c.Writer, c.Request)
	})

	// if Allow DirectoryIndex
	//r.Use(static.Serve("/", static.LocalFile("/tmp", true)))
	// set prefix
	//r.Use(static.Serve("/static", static.LocalFile("/tmp", true)))

	r.Use(static.Serve("/", static.LocalFile("../public", false)))

	r.Run("localhost:12312")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// websocket handler use gorilla/websocket
func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

// static file handler
func StaticServer(w http.ResponseWriter, req *http.Request) {
	file := dir + req.URL.Path
	fi, err := os.Stat(file)
	if os.IsNotExist(err) {
		http.NotFound(w, req)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if fi.IsDir() {
		if req.URL.Path[len(req.URL.Path)-1] != '/' {
			http.Redirect(w, req, req.URL.Path+"/", 301)
			return
		}
		for _, index := range indexs {
			fi, err = os.Stat(file + index)
			if err != nil {
				continue
			}
			http.ServeFile(w, req, file+index)
			return
		}
		http.NotFound(w, req)
		return
	}
	http.ServeFile(w, req, file)
}
