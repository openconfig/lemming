FROM golang:1.18-bullseye

WORKDIR /lemming

COPY . /lemming
RUN go build -o /lemming/dataplane-bin ./dataplane/main 

CMD [ "/lemming/dataplane-bin" ]