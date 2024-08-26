FROM golang:1.18 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.io

COPY . /app
WORKDIR /app

RUN go get && go build -ldflags="-s -w" -installsuffix cgo

FROM debian:buster-slim

ENV TZ=Asia/Shanghai \
    LANG=C.UTF-8 \
    APP_DIR=/usr/local/fastsearch

COPY --from=builder /app/fastsearch ${APP_DIR}/fastsearch

WORKDIR ${APP_DIR}

RUN ln -snf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && chmod +x fastsearch

EXPOSE 5679

CMD ["./fastsearch"]
