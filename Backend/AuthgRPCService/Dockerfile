FROM golang:1.23-alpine AS build

WORKDIR /AuthService

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o app ./cmd/AuthService

FROM alpine:latest
WORKDIR /root/

COPY --from=build /AuthService/app /AuthService/app
COPY config/config.yaml ./config/config.yaml
COPY logs ./logs

EXPOSE 8080

CMD ["/AuthService/app"]