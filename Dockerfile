FROM docker.io/golang:1.18-bullseye AS build
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /out/lemming ./cmd/lemming 

FROM docker.io/debian:bullseye
RUN apt-get update && apt-get -y install iproute2
COPY --from=build /out/lemming /lemming/lemming