# syntax=docker/dockerfile:1

FROM golang:1.18

WORKDIR /bpw
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/bpw ./...

EXPOSE 8080

CMD ["bpw"]
