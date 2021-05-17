FROM golang:1.15-alpine

ENV CONF_FILE_PATH=conf/docker/app.ini GO111MODULE=on GOPROXY=https://goproxy.cn GOSUMDB=off

RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories

RUN mkdir -p $GOPATH/src/github.com/Walker-PI/iot-gateway

COPY . $GOPATH/src/github.com/Walker-PI/iot-gateway

WORKDIR $GOPATH/src/github.com/Walker-PI/iot-gateway

RUN sh build.sh

EXPOSE 9922

ENTRYPOINT ["./output/bootstrap.sh"]