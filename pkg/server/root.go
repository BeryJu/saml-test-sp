package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"beryju.io/saml-test-sp/pkg/helpers"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
)

type Server struct {
	m *samlsp.Middleware
	h *http.ServeMux
	l *log.Entry
	b string
}

func (s *Server) hello(w http.ResponseWriter, r *http.Request) {
	sa := samlsp.SessionFromContext(r.Context())
	if s == nil {
		http.Error(w, "No Session", http.StatusInternalServerError)
		return
	}
	sa, ok := sa.(samlsp.SessionWithAttributes)
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
	_, _ = w.Write(data)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = fmt.Fprint(w, "hello :)")
}

func (s *Server) logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.l.WithField("remoteAddr", r.RemoteAddr).WithField("method", r.Method).Info(r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	nameID := samlsp.AttributeFromContext(r.Context(), "urn:oasis:names:tc:SAML:attribute:subject-id")
	var binding *saml.Endpoint
	for _, desc := range s.m.ServiceProvider.IDPMetadata.IDPSSODescriptors {
		for _, slo := range desc.SingleLogoutServices {
			s.l.WithField("slo", slo.Binding).Info("found SLO binding")
			binding = &slo
			break
		}
	}
	if binding == nil {
		s.l.Warning("no SLO descriptors found, aborting")
		w.WriteHeader(400)
		return
	}

	switch binding.Binding {
	case saml.HTTPRedirectBinding:
		url, err := s.m.ServiceProvider.MakeRedirectLogoutRequest(nameID, s.b+"/health")
		if err != nil {
			s.l.WithError(err).Warning("failed to make redirect logout")
		}
		http.Redirect(w, r, url.String(), http.StatusFound)
	case saml.HTTPPostBinding:
		res, err := s.m.ServiceProvider.MakePostLogoutRequest(nameID, s.b+"/health")
		if err != nil {
			s.l.WithError(err).Warning("failed to make post logout")
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(res)
	default:
		http.Error(w, "invalid binding", 500)
	}
	err := s.m.Session.DeleteSession(w, r)
	if err != nil {
		s.l.WithError(err).Warning("failed to delete session")
	}
}

func RunServer() {
	config := helpers.LoadConfig()

	samlSP, err := samlsp.New(config)

	if err != nil {
		panic(err)
	}
	server := Server{
		m: samlSP,
		h: http.NewServeMux(),
		l: log.WithField("component", "server"),
	}
	server.h.Handle("/", samlSP.RequireAccount(http.HandlerFunc(server.hello)))
	server.h.Handle("/saml/logout", samlSP.RequireAccount(http.HandlerFunc(server.logout)))
	server.h.Handle("/saml/", samlSP)
	server.h.HandleFunc("/health", server.health)

	listen := helpers.Env("SP_BIND", "localhost:9009")
	server.b = listen
	server.l.Infof("Server listening on '%s'", listen)
	server.l.Infof("ACS URL is '%s'", samlSP.ServiceProvider.AcsURL.String())

	if _, set := os.LookupEnv("SP_SSL_CERT"); set {
		server.l.Info("SSL enabled")
		// SP_SSL_CERT set, so we run SSL mode
		err := http.ListenAndServeTLS(listen, os.Getenv("SP_SSL_CERT"), os.Getenv("SP_SSL_KEY"), server.logRequest(server.h))
		if err != nil {
			panic(err)
		}
	} else {
		err = http.ListenAndServe(listen, server.logRequest(server.h))
		if err != nil {
			panic(err)
		}
	}
}
