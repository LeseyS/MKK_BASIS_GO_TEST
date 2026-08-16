FROM golang:1.25.8-alpine3.22 AS build

WORKDIR /app

# Modules layer
COPY go.mod go.sum ./
RUN go mod download

# Build layer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -o /task_service ./cmd/app

FROM alpine:3.22 AS run

COPY --from=build /task_service /task_service
COPY --from=build /app/.env /.env

EXPOSE 8080

CMD ["/task_service"]