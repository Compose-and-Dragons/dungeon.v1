FROM --platform=$BUILDPLATFORM golang:1.25.1-alpine AS builder
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
COPY --from=builder /app/zephyr .
EXPOSE 9090
ENTRYPOINT ["./zephyr"]
