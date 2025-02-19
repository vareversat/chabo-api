# syntax=docker/dockerfile:1

FROM golang:1.24.0 as build
ARG VERSION
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .
RUN swag init -d ./internal/api/routers,./ -g main_router.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/chabo-api -ldflags="-X github.com/vareversat/chabo-api/internal/api/routers.version=$VERSION"

FROM alpine:3.20
COPY --from=build /app/chabo-api /usr/bin/local/chabo-api
EXPOSE 8080

CMD [ "/usr/bin/local/chabo-api" ]
