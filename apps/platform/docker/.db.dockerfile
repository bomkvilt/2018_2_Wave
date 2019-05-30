FROM ubuntu:18.10

ENV DEBIAN_FRONTEND noninteractive

# install psql
RUN apt-get update;                 \
    apt-get install -y              \
        sudo                        \
        wget                        \
        software-properties-common  \
        postgresql                  \
        postgresql-contrib;

# setup psql
USER postgres
RUN /etc/init.d/postgresql start;                                           \
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';";  \
    createdb -O docker docker;                                              \
    /etc/init.d/postgresql stop;

ENV PG_VERSION 10
RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PG_VERSION/main/pg_hba.conf;          \
    echo "listen_addresses='*'"       >> /etc/postgresql/$PG_VERSION/main/postgresql.conf;      \
    echo "synchronous_commit = off"   >> /etc/postgresql/$PG_VERSION/main/postgresql.conf;
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
EXPOSE 5432

# copy schemas
COPY sql/ sql/
RUN /etc/init.d/postgresql start;       \ 
    psql -d docker -f sql/apps.sql;     \
    psql -d docker -f sql/users.sql;    \
    /etc/init.d/postgresql stop;

CMD service postgresql start; \
        sleep infinity;
