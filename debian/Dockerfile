FROM golang:1.20

RUN apt-get update -y && \
    apt-get dist-upgrade -y && \
    apt-get install -y \
        dh-golang \
        devscripts

COPY make-deb.sh /usr/local/bin/make-deb

VOLUME [ "/src" ]

CMD [ "/usr/local/bin/make-deb" ]
