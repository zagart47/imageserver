FROM golang:1.19-alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/app/main.go
FROM scratch
COPY --from=builder /build/main /bin/main
ENTRYPOINT ["/bin/main"]
EXPOSE 80 81