FROM golang:1.10.3-stretch as build
RUN mkdir -p /src/app
WORKDIR /src/app
ENV GOPATH=/ GOBIN=/go/bin
RUN go get -u github.com/golang/dep/...
COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure -vendor-only
COPY . .
RUN export GOOS=linux GOARCH=amd64 CGO_ENABLED=1 && go build -o gfile

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir /app
COPY --from=build /src/app/gfile /app/gfile
ENV GIN_MODE=release
CMD [ "/app/gfile" ]