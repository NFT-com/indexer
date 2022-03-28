FROM golang:1.16-buster AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /jobs-api ./cmd/jobs-api
RUN go build -o /jobs-watcher ./cmd/jobs-watcher
RUN go build -o /parsing-dispatcher ./cmd/parsing-dispatcher
RUN go build -o /chain-bootstrapper ./cmd/chain-bootstrapper
RUN go build -o /chain-watcher ./cmd/chain-watcher

FROM gcr.io/distroless/base-debian10 as jobs-api
WORKDIR /
COPY --from=build /jobs-api /jobs-api
EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/jobs-api"]

FROM gcr.io/distroless/base-debian10 as chain-bootstrapper
WORKDIR /
COPY --from=build /chain-bootstrapper /chain-bootstrapper
USER nonroot:nonroot
ENTRYPOINT ["/chain-bootstrapper"]

FROM gcr.io/distroless/base-debian10 as chain-watcher
WORKDIR /
COPY --from=build /chain-watcher /chain-watcher
USER nonroot:nonroot
ENTRYPOINT ["/chain-watcher"]

FROM gcr.io/distroless/base-debian10 as jobs-watcher
WORKDIR /
COPY --from=build /jobs-watcher /jobs-watcher
USER nonroot:nonroot
ENTRYPOINT ["/jobs-watcher"]

FROM gcr.io/distroless/base-debian10 as parsing-dispatcher
WORKDIR /
COPY --from=build /parsing-dispatcher /parsing-dispatcher
USER nonroot:nonroot
ENTRYPOINT ["/parsing-dispatcher"]