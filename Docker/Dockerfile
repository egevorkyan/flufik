FROM registry.access.redhat.com/ubi9/ubi:latest

LABEL maintainer="egevorkyan@outlook.com"

ARG UNAME=flufik

WORKDIR /flufik

COPY flufik /bin/flufik
COPY service.sh .

RUN chmod +x /bin/flufik \
    && useradd -m -d /opt/flufik $UNAME
RUN chown $UNAME:$UNAME /opt/flufik/
RUN chown $UNAME:$UNAME service.sh

CMD ["./service.sh"]

USER $UNAME
