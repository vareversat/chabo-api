# syntax=docker/dockerfile:1

FROM golang:1.21 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .
RUN swag init -d ./internal/api,./ -g router.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/chabo-api

FROM alpine:3.14
COPY --from=build /app/chabo-api /usr/bin/local/chabo-api
EXPOSE 8080

ENTRYPOINT [ "/usr/bin/local/chabo-api" ]
