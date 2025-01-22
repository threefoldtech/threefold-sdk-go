FROM golang:1.21-alpine as builder

WORKDIR /app

COPY /node-registrar/go.mod /node-registrar/go.sum .

RUN go mod download

COPY /node-registrar .

RUN go build -o server ./cmds/registrar.go

FROM alpine:3.21

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
