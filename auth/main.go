package main

import (
	"auth/security"
	"auth/user"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func main() {
	uService := user.NewService()
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.NewLogger(),
	)
	r.Handle("/v1/auth", n.With(
		negroni.Wrap(userAuth(uService)),
	)).Methods("POST", "OPTIONS")

	r.Handle("/v1/validate-token", n.With(
		negroni.Wrap(validateToken()),
	)).Methods("POST", "OPTIONS")

	http.Handle("/", r)
	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":8081", //@TODO usar variável de ambiente
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func userAuth(uService user.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var param struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		err = uService.ValidateUser(param.Email, param.Password)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		var result struct {
			Token string `json:"token"`
		}
		result.Token, err = security.NewToken(param.Email)
		if err != nil {
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

func validateToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var param struct {
			Token string `json:"token"`
		}
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		t, err := security.ParseToken(param.Token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tData, err := security.GetClaims(t)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var result struct {
			Email string `json:"email"`
		}
		result.Email = tData["email"].(string)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		return
	})
}
