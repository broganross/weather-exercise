FROM golang:1.22-alpine3.19 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o ./weather-exercise ./cmd/server/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=build /src/weather-exercise ./
EXPOSE 80
ENTRYPOINT ["/app/weather-exercise"]