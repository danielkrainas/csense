FROM frolvlad/alpine-glibc

COPY ./dist /bin/csense

WORKDIR /

CMD /bin/csense
