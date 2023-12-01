FROM docker.io/library/debian:12-slim
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean
COPY saml-test-sp /
EXPOSE 9009
WORKDIR /web-root
ENV SP_BIND=0.0.0.0:9009
HEALTHCHECK --interval=5s --start-period=1s CMD [ "wget", "--spider", "http://localhost:9009/health" ]
ENTRYPOINT [ "/saml-test-sp" ]
