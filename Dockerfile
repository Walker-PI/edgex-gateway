FROM golang:1.15

ENV CONF_FILE_PATH=config/docker/app.ini GO111MODULE=on GOPROXY=https://goproxy.cn GOSUMDB=off

RUN mkdir -p $GOPATH/src/github.com/nju-iot/edgex_admin

COPY . $GOPATH/src/github.com/nju-iot/edgex_admin/

COPY ./wait-for /usr/local/bin/

WORKDIR $GOPATH/src/github.com/nju-iot/edgex_admin

RUN sh build.sh

EXPOSE 9922

ENTRYPOINT ["./output/bootstrap.sh"]