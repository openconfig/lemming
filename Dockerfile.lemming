FROM docker.io/golang:1.18-bullseye AS build
WORKDIR /build
COPY go.* ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /out/lemming ./cmd/lemming 

FROM docker.io/debian:bullseye
COPY --from=build /out/lemming /lemming/lemming