# syntax=docker/dockerfile:1
FROM golang:1.22.5

WORKDIR /rest-arch-training

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /server

EXPOSE 8080

CMD ["/server"]