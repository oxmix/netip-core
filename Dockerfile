FROM docker.io/library/golang:1.23-alpine AS builder
ARG GO_PRIVATE
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o core .

FROM docker.io/library/alpine:3.20
LABEL description="https://cloudnetip.com/wiki"
ARG VERSION
ENV VERSION=$VERSION
ARG VERSION_HASH
ENV VERSION_HASH=$VERSION_HASH
# who logged
RUN apk --no-cache add dbus elogind
# benchmarks
RUN apk --no-cache add sysbench fio speedtest-cli
# device metrics
RUN apk --no-cache add gcompat mdadm smartmontools
COPY --from=builder /app/core .
ENTRYPOINT ["./core"]
