package main

import (
	"fmt"
	"github.com/johan-scriptdrift/data-star-test/routes"
	"github.com/johan-scriptdrift/data-star-test/views"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello World")

	dsRouter := routes.NewDsRouter()
	router := routes.NewRouter(dsRouter)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := views.Index().Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
