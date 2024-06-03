FROM golang:1.20.7 AS build-stage

WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o /IPFS-CLUSTER-MANAGER ./cmd

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /IPFS-CLUSTER-MANAGER /app/.env ./

EXPOSE 8090

USER nonroot:nonroot

CMD ["/IPFS-CLUSTER-MANAGER"]
