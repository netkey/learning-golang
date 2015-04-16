// web log to influxdb plugin for gin
package influxloger

import (
	"log"
	//"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/influxdb/influxdb/client"
)


func Logger(client *influxdb.Client) gin.HandlerFunc {
	//return func(res http.ResponseWriter, req *http.Request, c martini.Context, log *log.Logger) {
	return func(c *gin.Context) {

		start := time.Now()
		//log.Printf("Started %s %s", req.Method, req.URL.Path)

		//rw := res.(martini.ResponseWriter)

		c.Next()

		t := time.Since(start)
		//log.Printf("Completed %v %s in %v\n", rw.Status(), http.StatusText(rw.Status()), t)
		status := c.Writer.Status()


		if client != nil {
			s := &influxdb.Series{
				Name:    "resp_time",
				Columns: []string{"duration", "code", "url", "method"},
				Points: [][]interface{}{
					[]interface{}{int64(t / time.Millisecond), status, c.ClientIP, req.Method},
				},
			}
			err := client.WriteSeries([]*influxdb.Series{s})
			if err != nil {
				log.Println(err)
			}
		}


	}
}


/*

r := gin.Default()


conf := &influxdb.ClientConfig{
    Host:     "lecaire.nobugware.com:8086",
    Username: "root",
    Password: "totoin",
    Database: "four",
}
client, err := influxdb.NewClient(conf)
if err != nil {
    log.Fatal(err)
}
r.Use(influxlogger.Logger(client))


 */


