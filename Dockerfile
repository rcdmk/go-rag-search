FROM golang:1.24 as builder

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -o main main.go

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /go/src/app/main .
ENTRYPOINT ["./main"]