FROM golang:alpine as builder
RUN apk --no-cache add git
WORKDIR /app
COPY . .
RUN GOOS=linux go build -ldflags "-s -w" -o bin/server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir -p /app/web
WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/web /app/web
CMD ["./server"]
