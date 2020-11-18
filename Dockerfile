FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/cmd/logger
COPY src $GOPATH/src
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/logger
FROM scratch
COPY --from=builder /go/bin/logger /go/bin/logger
USER 9999:9999
ENTRYPOINT ["/go/bin/logger"]
