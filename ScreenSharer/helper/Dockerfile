FROM golang:1.20rc2-alpine3.17 AS temp
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src
COPY . /go/src
ENV GO111MODULE=on
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/
RUN cd /go/src &&  \
    go mod download
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories &&  \
    apk add gcc g++ libffi-dev make zlib-dev libcec-dev libtool
RUN go build -o main ./main.go

FROM alpine:3.16.3
ENV HOST=127.0.0.1
ENV PORT=8080
WORKDIR /app
COPY --from=temp /go/src/main /app/app
COPY --from=temp /go/src/html /app/html
COPY --from=temp /go/src/image/notfound.png /app/image/notfound.png
VOLUME /app/image
VOLUME /app/data
EXPOSE 8080
ENTRYPOINT ["./app"]