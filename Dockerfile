FROM frolvlad/alpine-glibc
MAINTAINER Danny Krainas <me@danielkrainas.com>

ENV CSENSE_CONFIG_PATH /etc/csense.default.yml

COPY ./dist /bin/csense
COPY ./config.default.yml /etc/csense.default.yml

ENTRYPOINT ["/bin/csense"]
