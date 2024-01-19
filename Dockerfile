FROM golang:1.21-alpine as go_builder

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/main.go

FROM scratch

COPY --from=go_builder /app/main /app/main

EXPOSE 3000
ENTRYPOINT ["/app/main"]
CMD ["serve"]