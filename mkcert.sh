#!/bin/bash
openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
  -keyout saml-sp.key -out saml-sp.pem -subj "/CN=localhost"
export SP_SSL_CERT=./saml-sp.pem
export SP_SSL_KEY=./saml-sp.key
