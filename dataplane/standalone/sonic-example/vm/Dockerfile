FROM ubuntu:20.04

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -qy \
 && apt-get upgrade -qy \
 && apt-get install -y \
    bridge-utils \
    iproute2 \
    socat \
    qemu-kvm \
    tcpdump \
    inetutils-ping \
    dnsutils \
    iptables \
    vim \
    wget \
    uml-utilities \
    isc-dhcp-client \
    dnsmasq \
 && rm -rf /var/lib/apt/lists/*

RUN wget https://go.dev/dl/go1.20.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz
ENV PATH="${PATH}:/usr/local/go/bin"

COPY vm.img /vm.img
COPY go.* /
RUN go mod download
COPY *.go /
COPY tc-tap.sh /

EXPOSE 22
CMD [ "go" "run" "/main.go" ] 
