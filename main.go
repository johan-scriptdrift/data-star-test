package main

import (
	"context"
	"fmt"
	"github.com/johan-scriptdrift/data-star-test/routes"
	"github.com/johan-scriptdrift/data-star-test/sql"
	"github.com/johan-scriptdrift/data-star-test/sql/zz"
	"github.com/johan-scriptdrift/data-star-test/views"
	zl "github.com/rs/zerolog/log"
	"github.com/starfederation/datastar-go/datastar"
	"log"
	"net/http"
	"time"
	"zombiezen.com/go/sqlite"
)

func main() {
	fmt.Println("Hello World")

	db, err := sql.SetupDB(context.Background(), "data", true)
	if err != nil {
		zl.Fatal().Err(err).Msg("Failed to setup database")
	}

	defer db.Close()

	dsRouter := routes.NewDsRouter()
	router := routes.NewRouter(dsRouter)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := views.Index([]zz.LocationModel{}).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	dsRouter.Get("/update-greeting", func(w http.ResponseWriter, r *http.Request) {
		err := db.ReadTX(r.Context(), func(tx *sqlite.Conn) error {
			u, err := zz.OnceReadByIDUser(tx, 1)
			if err != nil {
				return err
			}
			log.Println(u.FirstName)
			w.WriteHeader(200)
			sse := datastar.NewSSE(w, r)
			if err := sse.PatchElementTempl(views.UpdateGreeting(u.FirstName)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	dsRouter.Get("/locations", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		sse := datastar.NewSSE(w, r)
		if err := sse.PatchElementTempl(views.UpdateGreeting("SSE")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := db.ReadTX(r.Context(), func(tx *sqlite.Conn) error {
			locations, err := zz.OnceReadAllLocations(tx)
			if err != nil {
				return err
			}

			if err := sse.PatchElementTempl(views.Index(locations)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			streamedLocations := make([]zz.LocationModel, 0)
			for _, l := range locations {
				streamedLocations = append(streamedLocations, l)

				if err := sse.PatchElementTempl(views.Index(streamedLocations)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				time.Sleep(1 * time.Second)
			}

			return nil
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
