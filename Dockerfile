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
CMD ["./rd"]
