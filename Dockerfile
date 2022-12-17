FROM debian:bullseye as build

RUN apt update && apt install -y git golang ca-certificates

WORKDIR /src

ENV GOBIN=/src/bin \
    GOPATH=/src/go \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .

RUN go mod download

ADD . /src/
RUN go install ./...

FROM scratch

COPY --from=build /src/bin/server /server
COPY --from=build /src/bin/client /client
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/server"]
