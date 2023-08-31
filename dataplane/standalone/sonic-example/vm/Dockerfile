FROM ubuntu:20.04

RUN apt-get update -y \
 && apt-get install -y \
    bridge-utils \
    iproute2 \
    qemu-kvm \
    tcpdump \
    inetutils-ping \
    vim \
    wget

RUN wget https://go.dev/dl/go1.20.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz
ENV PATH="${PATH}:/usr/local/go/bin"

COPY vm.img /vm.img
COPY go.* /
RUN go mod download
COPY *.go /
COPY run.sh /
COPY tc-tap.sh /

EXPOSE 22
CMD [ "go" "run" "/main.go" ] 
