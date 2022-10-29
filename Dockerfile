FROM golang:1.19 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'github.com/egevorkyan/flufik/core.Version=1.5-2'" -o flufik cmd/main/main.go


FROM registry.access.redhat.com/ubi9/ubi:latest AS production

LABEL maintainer="egevorkyan@outlook.com"

ARG UNAME=flufik

WORKDIR /flufik

COPY --from=builder /app/flufik /bin/flufik
COPY service.sh .

RUN chmod +x /bin/flufik \
    && useradd -m -d /opt/flufik $UNAME
RUN chown $UNAME:$UNAME /opt/flufik/
RUN chown $UNAME:$UNAME service.sh

CMD ["./service.sh"]

USER $UNAME