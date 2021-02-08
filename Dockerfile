# Build stage
FROM golang AS build-env
ADD . /src/go-grpc-laughing-broccoli
ENV CGO_ENABLED=0
RUN cd /src/go-grpc-laughing-broccoli && go build -o /app

# Production stage
FROM scratch
COPY --from=build-env /app /

ENTRYPOINT ["/app"]
