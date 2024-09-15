# syntax=docker/dockerfile:1
FROM golang:1.22.5

WORKDIR /rest-arch-training

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /server ./cmd/rest-arch-training

EXPOSE 8080

CMD ["/server", "/rest-arch-training/config/prod.yaml"]