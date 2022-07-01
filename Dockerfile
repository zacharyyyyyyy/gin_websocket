FROM alpine:latest

ENV GOPROXY https://goproxy.cn

ADD . /go/gin_websocket
# Timezone
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' > /etc/timezone

WORKDIR /go/gin_websocket

CMD ["/go/gin_websocket/main"]