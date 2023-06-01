# fn

Using GIN in AWS Lambda

```go
package main

import (
	"github.com/StevenZack/fn"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := fn.Default()
	r.GET("/", func(c *gin.Context) {
		m := gin.H{}
		for k := range c.Request.Header {
			m[k] = c.GetHeader(k)
		}
		c.JSON(200, m)
	})
	r.Run()
}

```
