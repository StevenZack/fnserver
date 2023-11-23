package fn

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

func Run(serveMux http.Handler) error {
	lambda.Start(func(ctx context.Context, r request) (*response, error) {
		// convert to gin.Context
		req, e := r.ToHttpRequest()
		if e != nil {
			log.Println(e)
			return jsonResponse(500, map[string]any{
				"message": e.Error(),
			})
		}
		w := newResponse()

		serveMux.ServeHTTP(w, req)
		w.Wrap()
		return w, nil
	})
	return nil
}
