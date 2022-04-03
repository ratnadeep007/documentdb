package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func newServer(database string, port string) (*server, error) {
	s := server{db: nil, port: port}
	var err error
	s.db, err = pebble.Open(database, &pebble.Options{})
	return &s, err
}

func main() {
	databaseName, port := parseCommandLine(os.Args)

	s, err := newServer(databaseName, port)
	if err != nil {
		log.Fatal(err)
	}
	defer s.db.Close()

	router := httprouter.New()
	router.GET("/status", s.status)
	router.POST("/docs", s.addDocument)
	router.GET("/search", s.searchDocuments)
	router.GET("/docs/:id", s.getDocument)

	handler := cors.Default().Handler(router)
	log.Println("Listening on " + s.port)
	log.Fatal(http.ListenAndServe(":"+s.port, handler))
}
