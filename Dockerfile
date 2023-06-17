FROM registry.cn-hangzhou.aliyuncs.com/wl-common/golang:alpine AS builder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct


WORKDIR /build/zero

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY etc/taskschedule-api-${ACTIVE}.yaml /app/etc/taskschedule-api.yaml
COPY etc/config-${ACTIVE}.toml /app/etc/config.toml
RUN go build -ldflags="-s -w" -o /app/dialogue-manager taskschedule/taskschedule.go

FROM registry.cn-hangzhou.aliyuncs.com/wl-common/alpine

EXPOSE 8080
WORKDIR /app
COPY --from=builder /app/dialogue-manager /app/dialogue-manager
COPY --from=builder /app/etc /app/etc

CMD ["./dialogue-manager", "-f", "etc/taskschedule-api.yaml"]
