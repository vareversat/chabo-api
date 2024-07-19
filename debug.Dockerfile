# syntax=docker/dockerfile:1

FROM golang:1.22.4 as build
ARG VERSION
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .
RUN swag init -d ./internal/api/routers,./ -g main_router.go
# Debug mode
RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/chabo-api -gcflags "all=-N -l" -ldflags="-X github.com/vareversat/chabo-api/internal/api/routers.version=$VERSION"

CMD [ "/go/bin/dlv", "--listen=:4000", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/app/chabo-api" ]
