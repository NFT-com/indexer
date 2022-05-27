FROM golang:1.17-buster AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /watcher ./cmd/jobs-watcher

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /watcher /watcher
USER nonroot:nonroot
ENTRYPOINT ["/watcher"]
