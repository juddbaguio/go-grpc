FROM golang:1.19.0-bullseye AS builder

WORKDIR /app
RUN ls
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/auth ./cmd/auth


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/auth .
EXPOSE 3000
CMD [ "./auth" ]

