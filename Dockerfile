FROM golang:1.12
ENV GO111MODULE on
WORKDIR /go/src/github.com/goraft/raftd
ADD . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o raftd

FROM alpine:3.10
COPY --from=0 /go/src/github.com/goraft/raftd/raftd /
ENTRYPOINT ["/raftd"]