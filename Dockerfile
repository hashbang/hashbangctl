FROM debian:buster as build

RUN apt update && apt install -y git golang

ADD . /src/

RUN cd /src/ && \
    GOBIN=/src/bin \
    GOPATH=/src/go \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go install ./...

FROM scratch

COPY --from=build /src/bin/server /server
COPY --from=build /src/bin/signup /signup

CMD ["/server"]
