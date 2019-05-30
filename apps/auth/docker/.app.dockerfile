FROM ubuntu:18.10

ENV DEBIAN_FRONTEND noninteractive

# install golang
RUN apt-get update;                                     \
    apt-get install -y                                  \
        software-properties-common;                     \
    add-apt-repository ppa:longsleep/golang-backports;  \
    add-apt-repository ppa:apt-fast/stable;             \
    apt-get install -y                                  \
        apt-fast;                                       \
    apt-fast install -y                                 \
        golang-go;

# ports
EXPOSE 50000

USER root
ARG  ROOT="."
WORKDIR /home/docker/build
COPY . .

CMD go run $ROOT/cmd/;
