FROM golang:1.13-stretch AS builder

WORKDIR /build

COPY . .

RUN make release

FROM golang:1.13-stretch

WORKDIR /app

COPY --from=builder /build/release/wejay /app

ENV \
    # UDP_SERVER=localhost:8090 \
    UDP_SERVER= \
    # ADDR= \
    ADDR= \
    PORT=8080 \
    HOST=0.0.0.0


EXPOSE 8080/tcp

RUN adduser --disabled-password --gecos '' wejay && \
    chmod -R g+rwX         /app && \
    chgrp -R wejay         /app && \
    chown -R wejay:wejay   /app

USER wejay

ENTRYPOINT [ "/app/bin" ]
CMD [  ]
