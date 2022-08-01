FROM golang:1.17-buster AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /api ./cmd/api

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /api /api
USER nonroot:nonroot
ENTRYPOINT ["/api"]
