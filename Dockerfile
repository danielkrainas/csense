FROM frolvlad/alpine-glibc

COPY ./dist /bin/csense

WORKDIR /

ENTRYPOINT ["/bin/csense"]
