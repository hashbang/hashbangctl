FROM debian:buster

RUN adduser admin && \
    apt-get clean && \
    apt-get update && \
    apt-get install -y --force-yes \
        bats \
        sudo \
        ssh &&\
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ADD . /home/admin/

RUN \
    printf "source /home/admin/test_helper.bash\nsetup" \
        >> /home/admin/.bashrc && \
    chown -R admin:admin /home/admin && \
    echo '%admin ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

USER admin

WORKDIR /home/admin

ENV LANG=C.UTF-8 \
    TZ=UTC \
    TERM=xterm-256color \
    USER="admin" \
    HOME="/home/admin"

CMD ["/usr/bin/bats", "/test/test.bats"]
