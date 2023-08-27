FROM golang:1.21 as builder

WORKDIR /go/src
COPY . .
RUN ls
RUN cd internal/post-forwarder && CGO_ENABLED=0 go build -o "/go/bin/pf" -ldflags="-w -s"

FROM alpine:3.18.2
RUN adduser -S -D -H -h /app appuser
USER appuser
WORKDIR /app
COPY --from=builder /go/bin/pf /app/

EXPOSE 4321
CMD /app/pf
