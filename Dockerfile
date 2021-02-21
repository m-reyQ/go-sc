FROM golang:1.10

RUN apt-get update && apt-get install -y curl xvfb xorg-dev libglu1-mesa-dev

WORKDIR /go/src/github.com/d4l3k/go-sct
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
