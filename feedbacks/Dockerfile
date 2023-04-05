FROM golang:1.20-alpine AS builder

WORKDIR /src

# Restore dependencies
COPY feedbacks/ .

RUN go mod tidy

# Build executable
RUN go build -o /src/feedbacks ./


FROM scratch
WORKDIR /src
COPY --from=builder /src/feedbacks ./
EXPOSE 8082
CMD ["/src/feedbacks"]
