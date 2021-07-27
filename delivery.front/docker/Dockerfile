# Build stage.
FROM golang:1.15-alpine3.12 as build

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /opt/faem/deliveryfront

# Import the code from the context.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 go build \
  -trimpath \
  -mod=vendor \
  -o bin/deliveryfront \
  ./cmd/

# Final stage: the running container.
FROM alpine:3.12 AS final

# Install runtime dependencies and set up a user.
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
  && addgroup -S faem -g 1337 \
  && adduser -S -G faem delivery -u 1337

# Import the compiled executable and config file from the first stage.
COPY --from=build /opt/faem/deliveryfront/bin/deliveryfront /usr/local/bin/deliveryfront
COPY --from=build /opt/faem/deliveryfront/deployment/config/deliveryfront.toml /etc/faem/deliveryfront.toml

USER 1337

ENV TZ Europe/Moscow

# Run the compiled binary and pass the config.
CMD ["/usr/local/bin/deliveryfront", "-config", "/etc/faem/deliveryfront.toml"]
