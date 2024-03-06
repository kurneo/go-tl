FROM postgres:15.3

MAINTAINER GiangNguyen <giangnguyen.neko.130@gmail.com>

COPY ./init.sh /docker-entrypoint-initdb.d/init.sh

RUN chmod +x /docker-entrypoint-initdb.d/init.sh

EXPOSE 3306
