FROM golang:1.19-alpine as builder

COPY .  /app/
WORKDIR /app/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o courses-file-service .

FROM alpine
WORKDIR /app/
COPY --from=builder /app/courses-file-service .
COPY config/*.yml ./config/
CMD [ "./courses-file-service" ]
