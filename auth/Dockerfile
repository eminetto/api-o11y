FROM golang:1.21-alpine AS builder

WORKDIR /src

# Restore dependencies
COPY auth/ .

RUN go mod tidy

# Build executable
RUN go build -o /src/auth ./

FROM scratch
WORKDIR /src
COPY --from=builder /src/auth ./
EXPOSE 8081
CMD ["/src/auth"]
