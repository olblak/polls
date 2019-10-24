### Build Environment
FROM golang:1.13 as build

EXPOSE 80

VOLUME /var/log/backend

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/polls .

#### Polls image

FROM alpine:latest

EXPOSE 8080

LABEL \
  maintainer="https://github.com/olblak"

RUN apk --no-cache add ca-certificates

COPY --from=0 /go/bin/polls /polls

ENTRYPOINT ["./polls"]
