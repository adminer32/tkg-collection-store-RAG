FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/tkg-rag .

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache wget \
    && adduser -D -H rag
COPY --from=builder /out/tkg-rag /app/tkg-rag
COPY data /app/data
USER rag
ENV RAG_ADDR=:8081
ENV RAG_DATA_DIR=/app/data
EXPOSE 8081
CMD ["/app/tkg-rag"]
