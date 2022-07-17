FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /docker-gs-ping ./cmd/main.go

EXPOSE 3000

ENTRYPOINT [ "/docker-gs-ping" ]