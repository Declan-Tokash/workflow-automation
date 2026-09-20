FROM python:3.12-slim

# Install Git and basic utilities
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        git \
        ca-certificates \
        curl && \
    rm -rf /var/lib/apt/lists/*

# Create workspace
RUN mkdir -p /workspace

WORKDIR /workspace

# Make Python output appear immediately in container logs
ENV PYTHONUNBUFFERED=1

CMD ["sleep", "3600"]