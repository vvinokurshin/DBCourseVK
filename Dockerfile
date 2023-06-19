FROM golang:latest AS build

COPY . /server/

WORKDIR /server/

RUN CGO_ENABLED=0 go build cmd/main.go

FROM ubuntu:20.04
COPY . .

ENV TZ=Russia/Moscow
RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER valeriy WITH SUPERUSER PASSWORD 'valeriy_pw';" &&\
    createdb -O valeriy db_forum &&\
    psql -f db/db.sql -d db_forum &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

USER root
COPY --from=build /server/main .

EXPOSE 5000

CMD service postgresql start && ./main
