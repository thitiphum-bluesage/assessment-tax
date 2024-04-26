FROM golang:1.21-alpine as builder

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o main .

FROM alpine:latest

COPY --from=builder /go/src/app/main /main

ENV DATABASE_URL=postgres://godou:1111@postgres:5432/ktaxes?sslmode=disable
ENV PORT=8080
ENV ADMIN_USERNAME=adminTax
ENV ADMIN_PASSWORD=admin!

EXPOSE 8080

CMD ["./main"]
