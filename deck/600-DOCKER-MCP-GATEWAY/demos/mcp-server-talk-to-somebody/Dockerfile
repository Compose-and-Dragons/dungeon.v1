FROM --platform=$BUILDPLATFORM golang:1.24.4-alpine AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

RUN <<EOF
go mod tidy 
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build
EOF

FROM scratch
FROM alpine:latest
# for health check
RUN apk --no-cache add ca-certificates wget 
WORKDIR /app
COPY --from=builder /app/mcp-dd .
EXPOSE 9090
ENTRYPOINT ["./mcp-dd"]

# docker build --platform linux/arm64 -t mcp-dd:demo .
# NOTE: Inspector: docker run --rm -i mcp-dd:demo 
