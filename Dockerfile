FROM docker.io/library/alpine:3.18.0
RUN apk add --no-cache ca-certificates
COPY saml-test-sp /
EXPOSE 9009
WORKDIR /web-root
ENV SP_BIND=0.0.0.0:9009
HEALTHCHECK --interval=5s --start-period=1s CMD [ "wget", "--spider", "http://localhost:9009/health" ]
ENTRYPOINT [ "/saml-test-sp" ]
