FROM golang:1.25.2-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o analyzer main.go
ENTRYPOINT ["./analyzer"]