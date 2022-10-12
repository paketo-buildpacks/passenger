FROM ubuntu:bionic

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get -y install curl binutils build-essential libssl-dev

COPY entrypoint /entrypoint

ENTRYPOINT ["/entrypoint"]