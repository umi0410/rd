FROM golang:1.19-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o rd

FROM alpine
WORKDIR /app
COPY --from=builder /app/rd .
COPY ./config/default.yaml /app/config/default.yaml
COPY ./assets/favicon.ico /app/assets/favicon.ico
ENV RD_CONFIG_NAME=default
ENTRYPOINT ["./rd"]
