FROM golang:1.16-buster AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /indexer ./cmd/indexer

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /indexer /indexer
USER nonroot:nonroot
ENTRYPOINT ["/indexer","wss://mainnet.infura.io/ws/v3/668ea84fefe146b295f3c3714839a223","ethereum","mainnet"]