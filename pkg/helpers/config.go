package helpers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
)

func Env(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}

func LoadConfig() samlsp.Options {
	samlOptions := samlsp.Options{
		AllowIDPInitiated: true,
	}

	samlOptions.EntityID = Env("SP_ENTITY_ID", "saml-test-sp")

	metadataURL, metadataURLexists := os.LookupEnv("SP_METADATA_URL")
	if metadataURLexists {
		log.Debugf("Will attempt to load metadata from %s", metadataURL)
		idpMetadataURL, err := url.Parse(metadataURL)
		if err != nil {
			panic(err)
		}
		desc, err := samlsp.FetchMetadata(context.TODO(), http.DefaultClient, *idpMetadataURL)
		if err != nil {
			panic(err)
		}
		samlOptions.IDPMetadata = desc
	} else {
		ssoURL := Env("SP_SSO_URL", "")
		binding := Env("SP_SSO_BINDING", saml.HTTPPostBinding)
		samlOptions.IDPMetadata = &saml.EntityDescriptor{
			EntityID: samlOptions.EntityID,
			IDPSSODescriptors: []saml.IDPSSODescriptor{
				{
					SingleSignOnServices: []saml.Endpoint{
						{
							Binding:  binding,
							Location: ssoURL,
						},
					},
				},
			},
		}
		if signingCert := Env("SP_SIGNING_CERT", ""); signingCert != "" {
			samlOptions.IDPMetadata.IDPSSODescriptors[0].KeyDescriptors = []saml.KeyDescriptor{
				{
					Use: "singing",
					KeyInfo: saml.KeyInfo{
						X509Data: saml.X509Data{
							X509Certificates: []saml.X509Certificate{
								{
									Data: signingCert,
								},
							},
						},
					},
				},
			}
		}
	}

	defaultURL := "http://localhost:9009"
	if _, ok := os.LookupEnv("SP_SSL_CERT"); ok {
		defaultURL = "https://localhost:9009"
	}
	rootURL := Env("SP_ROOT_URL", defaultURL)
	url, err := url.Parse(rootURL)
	if err != nil {
		panic(err)
	}
	samlOptions.URL = *url

	priv, pub := Generate(fmt.Sprintf("localhost,%s", url.Hostname()))
	samlOptions.Key = priv
	samlOptions.Certificate = pub
	if sign := Env("SP_SIGN_REQUESTS", "false"); strings.ToLower(sign) == "true" {
		samlOptions.Key = LoadRSAKey(os.Getenv("SP_SSL_KEY"))
		samlOptions.Certificate = LoadCertificate(os.Getenv("SP_SSL_CERT"))
		samlOptions.SignRequest = true
		log.Debug("Signing requests")
	}
	log.Debugf("Configuration Optons: %+v", samlOptions)
	return samlOptions
}
