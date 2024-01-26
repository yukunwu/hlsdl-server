FROM debian:stable
#RUN cp /etc/apt/sources.list /etc/apt/sources.list_bk
#RUN sed -i 's/deb.debian.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list
#RUN sed -i 's|security.debian.org/debian-security|mirrors.tuna.tsinghua.edu.cn/debian-security|g' /etc/apt/sources.list
RUN apt-get update -y && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone && mkdir -p /app
RUN apt-get -qq install -y --no-install-recommends ca-certificates curl

COPY ./entrypoint.sh /app/entrypoint.sh
COPY ./bin/hlsdl_server /app/hlsdl_server
WORKDIR /app

RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]