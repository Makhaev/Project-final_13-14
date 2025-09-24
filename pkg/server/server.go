package server

import (
	"fmt"
	"net/http"

	"github.com/Makhaev/projectname/pkg/api"
)

func Run() error {
	api.Init()
	port := 7540
	http.Handle("/", http.FileServer(http.Dir("web")))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
