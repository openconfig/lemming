FROM debian:bullseye
RUN apt-get update && apt-get install -y --no-install-recommends build-essential=12.9 clang=1:11.0-51+nmu5 \
    wget=1.21-1+deb11u1 protobuf-compiler=3.12.4-1+deb11u1 protobuf-compiler-grpc=1.30.2-3 \
    libgrpc++-dev=1.30.2-3 libgoogle-glog-dev=0.4.0-4 git=1:2.30.2-1+deb11u2 openjdk-11-jdk=11.0.22+7-1~deb11u1 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
RUN wget -q https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
RUN rm go1.21.6.linux-amd64.tar.gz
ENV PATH="$PATH:/usr/local/go/bin"
ENV JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64
ENV PATH="$PATH:$JAVA_HOME/bin"