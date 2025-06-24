FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod tidy
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /server .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /server /server

EXPOSE 5525

CMD ["./server"]