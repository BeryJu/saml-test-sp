package helpers

import (
	"fmt"
	"net/url"
	"os"

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
		Logger: log.WithField("component", "saml-lib"),
	}

	samlOptions.EntityID = Env("SP_ENTITY_ID", "saml-test-sp")

	metadataURL, metadataURLexists := os.LookupEnv("SP_METADATA_URL")
	if metadataURLexists {
		log.Debugf("Will attempt to load metadata from %s", metadataURL)
		idpMetadataURL, err := url.Parse(metadataURL)
		if err != nil {
			panic(err)
		}
		samlOptions.IDPMetadataURL = idpMetadataURL
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
	}

	rootURL := Env("SP_ROOT_URL", "http://localhost:9009")
	url, err := url.Parse(rootURL)
	if err != nil {
		panic(err)
	}
	samlOptions.URL = *url

	priv, pub := Generate(fmt.Sprintf("localhost,%s", url.Hostname()))
	samlOptions.Key = priv
	samlOptions.Certificate = pub
	log.Debugf("Configuration Optons: %+v", samlOptions)
	return samlOptions
}
