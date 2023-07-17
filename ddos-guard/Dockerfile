FROM golang:1.20.1-alpine AS builder

COPY . /go/src/github.com/yukels

WORKDIR /go/src/github.com/yukels

RUN go mod download
RUN CGO_ENABLED=0 go install -ldflags "-s -w" github.com/yukels/ddos-guard

FROM alpine:latest

RUN apk update && apk --no-cache add ca-certificates

RUN mkdir /app
RUN addgroup -S --gid 1122 ddos_group && adduser -S ddos_user -G ddos_group --uid 1122
RUN chown -R ddos_user.ddos_group /app
USER ddos_user

WORKDIR /app
COPY --from=builder /go/bin/ddos-guard /app/

EXPOSE 8081
ENTRYPOINT ["/app/ddos-guard"]
