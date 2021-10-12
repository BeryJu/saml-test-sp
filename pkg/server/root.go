package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"beryju.org/saml-test-sp/pkg/helpers"
	"github.com/crewjam/saml/samlsp"
)

func hello(w http.ResponseWriter, r *http.Request) {
	s := samlsp.SessionFromContext(r.Context())
	if s == nil {
		http.Error(w, "No Session", http.StatusInternalServerError)
		return
	}
	sa, ok := s.(samlsp.SessionWithAttributes)
	if !ok {
		http.Error(w, "Session has no attributes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data, err := json.MarshalIndent(sa, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "hello :)")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("remoteAddr", r.RemoteAddr).WithField("method", r.Method).Info(r.URL)
		handler.ServeHTTP(w, r)
	})
}

func RunServer() {
	config := helpers.LoadConfig()

	samlSP, err := samlsp.New(config)

	if err != nil {
		panic(err)
	}
	http.Handle("/", samlSP.RequireAccount(http.HandlerFunc(hello)))
	http.Handle("/saml/", samlSP)
	http.HandleFunc("/health", health)

	listen := helpers.Env("SP_BIND", "localhost:9009")
	log.Infof("Server listening on '%s'", listen)
	log.Infof("ACS URL is '%s'", samlSP.ServiceProvider.AcsURL.String())

	if _, set := os.LookupEnv("SP_SSL_CERT"); set {
		// SP_SSL_CERT set, so we run SSL mode
		err := http.ListenAndServeTLS(listen, os.Getenv("SP_SSL_CERT"), os.Getenv("SP_SSL_KEY"), logRequest(http.DefaultServeMux))
		if err != nil {
			panic(err)
		}
	} else {
		err = http.ListenAndServe(listen, logRequest(http.DefaultServeMux))
		if err != nil {
			panic(err)
		}
	}
}
