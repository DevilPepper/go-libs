package request

import (
	"fmt"
	"net/http"
)

func BaseDomain(req *http.Request) string {
	scheme := "https"
	if req.TLS == nil {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s", scheme, req.Host)
}
