FROM ubuntu:bionic
WORKDIR /root/
RUN apt-get update && \
    apt-get install --no-install-recommends -y tcpdump

COPY envoy-als-server .

CMD ["./envoy-als-server"]