package logging

import (
	"bytes"
	"io/ioutil"
	"time"

	"encoding/json"
	"github.com/gin-gonic/gin"
)

// LogInfo is a base log information
type LogInfo struct {
	ClientIP    string        `json:"ip"`
	Date        string        `json:"date"`
	Method      string        `json:"method"`
	RequestURI  string        `json:"uri"`
	Referer     string        `json:"referer,omitempty"`
	HTTPVersion string        `json:"httpVersion"`
	Size        int           `json:"size"`
	Status      int           `json:"status"`
	UserAgent   string        `json:"userAgent"`
	Latency     time.Duration `json:"latency"`
}

// AccessLog is a log information of user access
type AccessLog struct {
	LogInfo
	Error error `json:"error,omitempty"`
}

// ActivityLog is a log information of user action
type ActivityLog struct {
	LogInfo
	RequestBody map[string]interface{} `json:"requestBody,omitempty"`
	Extra       interface{}            `json:"extra,omitempty"`
}

// GenerateLogInfo generates base log information
func GenerateLogInfo(c *gin.Context, start time.Time) LogInfo {
	return LogInfo{
		ClientIP:    c.ClientIP(),
		Date:        start.Format(time.RFC3339),
		Method:      c.Request.Method,
		RequestURI:  c.Request.URL.RequestURI(),
		Referer:     c.Request.Referer(),
		HTTPVersion: c.Request.Proto,
		Size:        c.Writer.Size(),
		Status:      c.Writer.Status(),
		UserAgent:   c.Request.UserAgent(),
		Latency:     time.Now().Sub(start),
	}
}

// ConvertToMapFromBody converts to a map from a request body
func ConvertToMapFromBody(c *gin.Context) (m map[string]interface{}, err error) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	if len(b) != 0 {
		err = json.Unmarshal(b, &m)
	}
	return
}
