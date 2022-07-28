FROM golang:1.17

WORKDIR /data

ADD ./src/ /data

RUN go env -w GO111MODULE=on &&\
    go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

EXPOSE 8080

CMD ["go", "run", "main.go"]


