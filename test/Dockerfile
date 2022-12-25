FROM debian:buster

RUN adduser admin && \
    apt-get clean && \
    apt-get update && \
    apt-get install -y --force-yes \
        bats \
        sudo \
        postgresql-client \
        curl \
	tmux \
        ssh &&\
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /home/test

ENV LANG=C.UTF-8 \
    TZ=UTC \
    TERM=xterm-256color \
    HOME="/home/test"

CMD ["/usr/bin/bats", "/test/test.bats"]
