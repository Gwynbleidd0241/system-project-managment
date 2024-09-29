FROM golang:1.23-alpine AS build

WORKDIR /MainHTTP

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o app ./cmd/mainHTTP

FROM alpine:latest
WORKDIR /root/

COPY --from=build /MainHTTP/app /MainHTTP/app
COPY config/config.yaml ./config/config.yaml
COPY logs ./logs

EXPOSE 8080

CMD ["/MainHTTP/app"]