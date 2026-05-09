package adapter

import (
	"net/http"

	"github.com/devlucas-java/luca-s3/internal/delivery/http/response"
)

type AdaptHandler func(w http.ResponseWriter, r *http.Request) error

func Adapt(h AdaptHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := h(w, r)

		if err != nil {
			response.ResponseError(w, err)
		}
	}
}
