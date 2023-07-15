package api

import (
	"fmt"
	"net/http"
)

func (s *StoreHub) healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world")
}
