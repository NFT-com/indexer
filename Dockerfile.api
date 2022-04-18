FROM golang:1.17-buster AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /api ./cmd/jobs-api

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /api /api
EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/api"]
