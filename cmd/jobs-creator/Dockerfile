FROM golang:1.17-buster AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /creator ./cmd/jobs-creator

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /creator /creator
USER nonroot:nonroot
ENTRYPOINT ["/creator"]
