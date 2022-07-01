FROM golang:1.17-buster AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /dispatcher ./cmd/completion-dispatcher

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /dispatcher /dispatcher
USER nonroot:nonroot
ENTRYPOINT ["/dispatcher"]
