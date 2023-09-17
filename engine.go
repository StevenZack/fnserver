package fn

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"
)

type Engine struct {
	*gin.Engine
}

func Default() *Engine {
	return &Engine{
		Engine: gin.Default(),
	}
}

func (s *Engine) Run(addr ...string) error {
	lambda.Start(func(ctx context.Context, r request) (*response, error) {
		// convert to gin.Context
		req, e := r.ToHttpRequest()
		if e != nil {
			log.Println(e)
			return jsonResponse(500, gin.H{
				"message": e.Error(),
			})
		}
		w := newResponse()

		s.ServeHTTP(w, req)
		w.Wrap()
		return w, nil
	})
	return nil
}
