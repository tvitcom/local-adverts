FROM golang:alpine AS container
RUN apk update && \
    apk add curl \
            git \
            bash \
            make \
            ca-certificates && \
    rm -rf /var/cache/apk/*

# install migrate which will be used by entrypoint.sh to perform DB migration
ARG MIGRATE_VERSION=4.7.1
ADD https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz /tmp
RUN tar -xzf /tmp/migrate.linux-amd64.tar.gz -C /usr/local/bin && mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate

WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.* ./
RUN go mod download
RUN go mod verify

# copy source files and container the binary
COPY . .
RUN make container


FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
RUN mkdir -p /var/log/app
WORKDIR /app/
COPY --from=container /usr/local/bin/migrate /usr/local/bin
COPY --from=container /app/migrations ./migrations/
COPY --from=container /app/server .
COPY --from=container /app/deployments/entrypoint.sh .
COPY --from=container /app/config/*.yml ./config/
RUN ls -la
ENTRYPOINT ["./entrypoint.sh"]
