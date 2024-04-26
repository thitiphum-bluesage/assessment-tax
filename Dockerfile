FROM golang:1.21-alpine as builder

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o main .

FROM alpine:latest  

COPY --from=builder /go/src/app/main /main

EXPOSE 8080

CMD ["./main"]
