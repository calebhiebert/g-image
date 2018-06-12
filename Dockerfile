FROM golang:1.10.3-stretch as build
RUN mkdir /src
COPY . /src/
RUN cd /src && go get -d -v ./... && go build -o gfile
CMD ["/src/gfile"]

# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# RUN mkdir /app
# COPY --from=build /src/gfile /app/
# CMD [ "/app/gfile" ]