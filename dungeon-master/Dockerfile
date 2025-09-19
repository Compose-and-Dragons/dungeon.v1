#FROM golang:1.24.0-alpine AS builder
FROM --platform=$BUILDPLATFORM golang:1.24.4-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
# COPY go.mod .
# COPY main.go .
COPY . .

RUN <<EOF
go mod tidy 
#go build
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build
EOF

#FROM scratch
FROM alpine:3.22.1
# for health check
# RUN apk --no-cache add ca-certificates wget 
WORKDIR /app
COPY --from=builder /app/dungeon-master .

# I use volume instead of copying the data
#COPY data ./data

#ENTRYPOINT ["./dungeon-master"]
