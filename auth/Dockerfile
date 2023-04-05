FROM golang:1.19.2-alpine AS builder
#RUN apk add build-base

WORKDIR /app/

# Restore dependencies
COPY * /app/

RUN ls -la /
RUN go mod download

# Build executable
RUN go build -o /app/auth ./

FROM scratch

WORKDIR /app/

COPY --from=builder /app/auth ./
EXPOSE 8081
CMD ["/app/auth"]
