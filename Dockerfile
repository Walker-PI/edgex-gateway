FROM golang:1.15

ENV CONF_FILE_PATH=conf/docker/app.ini 
ENV CGO_ENABLED=0 GOOS=linux GO111MODULE=on GOPROXY=https://goproxy.cn GOSUMDB=off
ENV TZ=Asia/Shanghai

RUN mkdir -p $GOPATH/src/github.com/Walker-PI/iot-gateway

COPY . $GOPATH/src/github.com/Walker-PI/iot-gateway

WORKDIR $GOPATH/src/github.com/Walker-PI/iot-gateway

RUN sh build.sh

EXPOSE 9922

ENTRYPOINT ["./output/bootstrap.sh"]