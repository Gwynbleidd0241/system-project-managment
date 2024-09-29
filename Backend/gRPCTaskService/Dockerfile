FROM golang:1.23-alpine AS build

WORKDIR /TaskService

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o app ./cmd/task-service

FROM alpine:latest
WORKDIR /root/

COPY --from=build /TaskService/app /TaskService/app
COPY ./config/config.yaml ./config/config.yaml
COPY ./logs ./logs

CMD ["/TaskService/app"]