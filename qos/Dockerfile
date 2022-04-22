FROM docker.io/alpine
COPY qos-server /
COPY crt.pem /run/secrets/tls/
COPY key.pem /run/secrets/tls/
RUN chmod +wr /run/secrets/tls/crt.pem
RUN chmod +wr /run/secrets/tls/key.pem

ENTRYPOINT ["/qos-server"]
