FROM golang:1.20.1-alpine AS builder

COPY . /go/src/github.com/yukels

WORKDIR /go/src/github.com/yukels

RUN go mod download
RUN CGO_ENABLED=0 go install -ldflags "-s -w" github.com/yukels/test/hello

FROM alpine:latest

RUN apk update && apk --no-cache add ca-certificates

RUN mkdir /app

WORKDIR /app
COPY --from=builder /go/bin/hello /app/

EXPOSE 8080
ENTRYPOINT ["/app/hello"]
