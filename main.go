package main

import (
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

func (apiC *apiConfig) handlerResetCount(w http.ResponseWriter, req *http.Request) {
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
	}

	servMux.Handle("/app/", http.StripPrefix("/app", apiC.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	servMux.HandleFunc("GET /api/healthz", handlerReadiness)
	servMux.HandleFunc("POST /api/validate_chirp", handlerValidateChip)

	servMux.HandleFunc("GET /admin/metrics", apiC.handlerCountRequests)
	servMux.HandleFunc("POST /admin/reset", apiC.handlerResetCount)

	err = server.ListenAndServe()
	fmt.Println(err)
}
