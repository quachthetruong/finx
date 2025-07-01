FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get --no-install-recommends install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN groupadd -r nonroot && useradd -r -g nonroot nonroot

# Copy the binary to the production image from the builder stage (from gitlab runner).
COPY ./server /app/server

# Use an unprivileged user.
USER nonroot

# Run the web service on container startup.
CMD ["/app/server"]
