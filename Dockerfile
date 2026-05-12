FROM golang:1.25 AS builder

ENV GOPROXY https://goproxy.cn/

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p /app/templates && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/doc-generation-mcp-server ./cmd/server

FROM alpine:3.22.4

WORKDIR /app

COPY --from=builder /out/doc-generation-mcp-server /app/doc-generation-mcp-server
COPY --from=builder /app/templates /app/templates

EXPOSE 9103

ENTRYPOINT ["/app/doc-generation-mcp-server"]
