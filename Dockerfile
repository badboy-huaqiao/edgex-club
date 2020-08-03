# Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
# SPDX-License-Identifier: GPL-2.0 

FROM golang:1.13-alpine AS builder

MAINTAINER Zhang Huaqiao <yhzhq1989@163.com>

RUN cp /etc/apk/repositories /etc/apk/repositories.bak
RUN echo "https://mirrors.ustc.edu.cn/alpine/v3.6/main" > /etc/apk/repositories
RUN echo "https://mirrors.ustc.edu.cn/alpine/v3.6/community" >> /etc/apk/repositories
RUN cat /etc/apk/repositories

WORKDIR /go/src/edgex-club/

ENV GOPROXY https://goproxy.io

RUN apk update && apk --no-cache add ca-certificates && apk add git make

COPY . .

RUN make build

FROM alpine:3.6

RUN cp /etc/apk/repositories /etc/apk/repositories.bak
RUN echo "https://mirrors.ustc.edu.cn/alpine/v3.6/main" > /etc/apk/repositories
RUN echo "https://mirrors.ustc.edu.cn/alpine/v3.6/community" >> /etc/apk/repositories
RUN cat /etc/apk/repositories

RUN apk update && apk --no-cache add ca-certificates

WORKDIR /edgex-club/
COPY --from=builder /go/src/edgex-club/cmd/edgex-club .

EXPOSE 443
EXPOSE 8080

ENTRYPOINT ["./edgex-club", "-confpath=res/docker/configuration.toml","-prod=true"]
