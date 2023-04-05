FROM golang:1.20-alpine AS builder

WORKDIR /src

# Restore dependencies
COPY votes/ .

RUN go mod tidy

# Build executable
RUN go build -o /src/votes ./

FROM scratch
WORKDIR /src
COPY --from=builder /src/votes ./
EXPOSE 8083
CMD ["/src/votes"]
