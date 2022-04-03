package main

import (
	"encoding/json"
	"net/http"

	"github.com/cockroachdb/pebble"
	uuid "github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (s server) status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsonResponse(w, map[string]any{
		"status": "ok",
	}, nil)
}

func (s server) addDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := json.NewDecoder(r.Body)
	var document map[string]any
	err := dec.Decode(&document)
	if err != nil {
		jsonResponse(w, nil, err)
	}

	id := uuid.New().String()

	bs, err := json.Marshal(document)
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}
	err = s.db.Set([]byte(id), bs, pebble.Sync)
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}
	jsonResponse(w, map[string]any{"id": id}, nil)
}

func (s server) searchDocuments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q, err := parseQuery(r.URL.Query().Get("q"))
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}
	var documents []map[string]any
	iter := s.db.NewIter(nil)
	for iter.First(); iter.Valid(); iter.Next() {
		var document map[string]any
		err = json.Unmarshal(iter.Value(), &document)
		if err != nil {
			jsonResponse(w, nil, err)
			return
		}
		if q.match(document) {
			documents = append(documents, map[string]any{
				"id":   string(iter.Key()),
				"body": document,
			})
		}
	}
	jsonResponse(w, map[string]any{"documents": documents, "count": len(documents)}, nil)
}

func (s server) getDocument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if id == "0" {
		documents := []map[string]any{}
		documentIter := s.db.NewIter(nil)
		for documentIsValid := documentIter.First(); documentIter.Valid(); documentIsValid = documentIter.Next() {
			if documentIsValid {
				document, err := s.getDocumentById(documentIter.Key())
				if err != nil {
					jsonResponse(w, nil, err)
					return
				}
				documents = append(documents, map[string]any{
					"id":   string(documentIter.Key()),
					"body": document,
				})
			}
		}
		doc := map[string]any{
			"values": documents,
		}
		jsonResponse(w, doc, nil)
		return
	}

	document, err := s.getDocumentById([]byte(id))
	if err != nil {
		jsonResponse(w, nil, err)
		return
	}

	jsonResponse(w, map[string]any{
		"document": document,
	}, nil)
}

func (s server) getDocumentById(id []byte) (map[string]any, error) {
	valBytes, closer, err := s.db.Get(id)
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	var document map[string]any
	err = json.Unmarshal(valBytes, &document)
	return document, err
}
