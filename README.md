# SAML-test-SP

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/beryju/saml-test-sp/ci-build.yml?branch=main&style=for-the-badge)

This is a small, golang-based SAML Service Provider, to be used in End-to-end or other testing. It uses the https://github.com/crewjam/saml Library for the actual SAML Logic.

SAML-test-SP supports IdP-initiated Login flows, *however* RelayState has to be empty for this to work.

This tool is full configured using environment variables.

## URLs

- `http://localhost:9009/health`: Healthcheck URL, used by the docker healtcheck.
- `http://localhost:9009/saml/acs`: SAML ACS URL, needed to configure your IdP.
- `http://localhost:9009/saml/metadata`: SAML Metadata URL, needed to configure your IdP.
- `http://localhost:9009/saml/logout`: SAML Logout URL.
- `http://localhost:9009/`: Test URL, redirects to SAML SSO URL.

## Configuration

- `SP_BIND`: Which address and port to bind to. Defaults to `0.0.0.0:9009`.
- `SP_ROOT_URL`: Root URL you're using to access the SP. Defaults to `http://localhost:9009`.
- `SP_ENTITY_ID`: SAML EntityID, defaults to `saml-test-sp`.
- `SP_METADATA_URL`: Optional URL that metadata is fetched from. The metadata is fetched on the first request to `/`.
- `SP_COOKIE_NAME`: Custom name for the session cookie. Defaults to `token`. Use this to avoid cookie name conflicts with other applications.
---
- `SP_SSO_URL`: If the metadata URL is not configured, use these options to configure it manually.
- `SP_SSO_BINDING`: Binding Type used for the IdP, defaults to POST. Allowed values: `urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST` and `urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect`
- `SP_SIGNING_CERT`: PEM-encoded Certificate used for signing, with the PEM Header and all newlines removed.
---
Optionally, if you want to use SSL, set these variables
- `SP_SSL_CERT`: Path to the SSL Certificate the server should use.
- `SP_SSL_KEY`: Path to the SSL Key the server should use.
- `SP_SIGN_REQUESTS`: Enable signing of requests.

Note: If you're manually setting `SP_ROOT_URL`, ensure that you prefix that URL with https.

## Running

This service is intended to run in a docker container

```
# beryju.org is a vanity URL for ghcr.io/beryju
docker pull beryju.io/saml-test-sp
docker run -d --rm \
    -p 9009:9009 \
    -e SP_ENTITY_ID=saml-test-sp \
    -e SP_SSO_URL=http://id.beryju.io/... \
    beryju.io/saml-test-sp
```

Or if you want to use docker-compose, use this in your `docker-compose.yaml`.

```yaml
version: '3.5'

services:
  saml-test-sp:
    image: beryju.io/saml-test-sp
    ports:
      - 9009:9009
    environment:
      SP_METADATA_URL: http://some.site.tld/saml/metadata
    # If you don't want SSL, cut here
      SP_SSL_CERT: /fullchain.pem
      SP_SSL_KEY: /privkey.pem
    volumes:
      - ./fullchain.pem:/fullchain.pem
      - ./privkey.pem:/privkey.pem
```
