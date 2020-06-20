# SAML-test-SP

![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/beryju/saml-test-sp?style=flat-square)
![Docker pulls](https://img.shields.io/docker/pulls/beryju/saml-test-sp.svg?style=flat-square)

This is a small, golang-based SAML Service Provider, to be used in End-to-end or other testing. It uses the https://github.com/crewjam/saml Library for the actual SAML Logic.

This tool is full configured using environment variables.

## URLs

- `http://localhost:9009/health`: Healthcheck URL, used by the docker healtcheck.
- `http://localhost:9009/saml/acs`: SAML ACS URL, needed to configure your IdP.
- `http://localhost:9009/saml/metadata`: SAML Metadata URL, needed to configure your IdP.
- `http://localhost:9009/`: Test URL, redirects to SAML SSO URL.

## Configuration

- `SP_ROOT_URL`: Root URL you're using to access the SP. Defaults to `http://localhost:9000`.
- `SP_ENTITY_ID`: SAML EntityID, defaults to `saml-test-sp`
- `SP_METADATA_URL`: Optional URL that metadata is fetched from. The metadata is fetched on the first request to `/`
---
- `SP_SSO_URL`: If the metadata URL is not configured, use these options to configure it manually.
- `SP_SSO_BINDING`: Binding Type used for the IdP, defaults to POST Redirect. Allowed values: `urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST` and `urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect`
- `SP_SIGNING_CERT`: PEM-encoded Certificate used for signing, with the PEM Header and all newlines removed.

## Running

This service is intended to run in a docker container

```
docker pull beryju/saml-test-sp
docke run -d --rm \
    -p 9009:9009 \
    -e SP_ENTITY_ID=saml-test-sp \
    -e SP_SSO_URL=http://id.beryju.org/... \
    beryju/saml-test-sp
```
