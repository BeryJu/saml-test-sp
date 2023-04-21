FROM golang:latest AS builder
WORKDIR $GOPATH/src/github.com/BeryJu/saml-test-sp
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -v -o /go/bin/saml-test-sp

FROM alpine
COPY --from=builder /go/bin/saml-test-sp /saml-test-sp
EXPOSE 9009
WORKDIR /web-root
ENV SP_BIND=0.0.0.0:9009
HEALTHCHECK --interval=5s --start-period=1s CMD [ "wget", "--spider", "http://localhost:9009/health" ]
CMD [ "/saml-test-sp" ]
ENTRYPOINT [ "/saml-test-sp" ]
