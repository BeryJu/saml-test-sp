package server

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/BeryJu/saml-test-sp-go/pkg/helpers"
	"github.com/crewjam/saml/samlsp"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", samlsp.AttributeFromContext(r.Context(), "cn"))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "hello :)")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func RunServer() {
	config := helpers.LoadConfig()
	config.CookieSameSite = http.SameSiteNoneMode

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
		err := http.ListenAndServeTLS(listen, os.Getenv("SP_SSL_CERT"), os.Getenv("SP_SSL_KEY"), nil)
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
