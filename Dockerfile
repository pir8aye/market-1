FROM golang:latest
ADD . /go/src/market
RUN go install market
ENTRYPOINT ["/go/bin/market"]
