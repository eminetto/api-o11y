package main

import (
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/eminetto/talk-microservices-go/pkg/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
	"votes/vote"
)

func main() {
	vService := vote.NewService()
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.NewLogger(),
	)
	r.Handle("/v1/vote", n.With(
		negroni.HandlerFunc(middleware.IsAuthenticated()),
		negroni.Wrap(storeVote(vService)),
	)).Methods("POST", "OPTIONS")


	http.Handle("/", r)
	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":8083",//@TODO usar variável de ambiente
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func storeVote(vService vote.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var v vote.Vote
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		v.Email = r.Header.Get("email")
		var result struct {
			ID uuid.UUID `json:"id"`
		}
		result.ID, err = vService.Store(v)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		return
	})
}