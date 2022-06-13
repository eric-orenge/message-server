FROM golang:1.17 as build-env
RUN mkdir -p /go/src/github.com/eric-orenge/message-server
WORKDIR /go/src/github.com/eric-orenge/message-server
COPY go.mod . 
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -installsuffix cgo -o /go/bin/message-server
EXPOSE 5000
ENTRYPOINT ["/go/bin/message-server"]