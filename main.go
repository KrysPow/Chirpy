package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Kryspow/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{Content-Type: text/plain; charset=utf-8, Status: OK}"))
}

func (apiC *apiConfig) handlerCountRequests(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	s := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, apiC.fileServerHits.Load())
	w.Write([]byte(s))
}

func (apiC *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if apiC.platform != "dev" {
		w.WriteHeader(403)
		return
	}
	apiC.dbQueries.DeleteUsers(context.Background())
	w.WriteHeader(http.StatusOK)
	apiC.fileServerHits = atomic.Int32{}
	s := fmt.Sprintf("Hits: %v", apiC.fileServerHits.Load())
	w.Write([]byte(s))
}

func (apiC *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiC.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not connect database: ", err)
	}

	servMux := http.NewServeMux()
	server := http.Server{
		Handler: servMux,
		Addr:    ":8080",
	}

	apiC := &apiConfig{
		fileServerHits: atomic.Int32{},
		dbQueries:      database.New(db),
		platform:       os.Getenv("PLATFORM"),
	}

	servMux.Handle("/app/", http.StripPrefix("/app", apiC.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	servMux.HandleFunc("GET /api/healthz", handlerReadiness)
	servMux.HandleFunc("POST /api/users", apiC.handlerUsers)
	servMux.HandleFunc("POST /api/chirps", apiC.handlerChirps)

	servMux.HandleFunc("GET /admin/metrics", apiC.handlerCountRequests)
	servMux.HandleFunc("POST /admin/reset", apiC.handlerReset)

	err = server.ListenAndServe()
	fmt.Println(err)
}
