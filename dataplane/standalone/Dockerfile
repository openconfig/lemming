FROM debian:bullseye
RUN apt-get update && apt-get install -y build-essential clang wget protobuf-compiler protobuf-compiler-grpc libgrpc++-dev libgoogle-glog-dev git openjdk-11-jdk
RUN wget -q https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
RUN rm go1.21.6.linux-amd64.tar.gz
ENV PATH="$PATH:/usr/local/go/bin"
ENV JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64
ENV PATH="$PATH:$JAVA_HOME/bin"