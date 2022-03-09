FROM --platform=$BUILDPLATFORM golang:alpine3.12 as builder

ARG TARGETOS TARGETARCH TARGETVARIANT
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH

# This construct stores custom local variables in a local environment file which can be
# sourced in later build steps. It is necessary because docker doesn't support environment
# variables being set to the output of a command.
ENV LOCALENV="/local_env"
# For an architecture of `arm/v7`, extract the `7` and store it in GOARM
RUN if [ "arm" = "$TARGETARCH" ]; then \
        echo "export GOARM=$(echo $TARGETVARIANT | sed 's/v\(.*\)/\1/g')" >>$LOCALENV ;\
    else \
        touch $LOCALENV ;\
    fi
WORKDIR /build

# at first download all dependencies, so this step can be omitted on next run
COPY go.mod .
COPY go.sum .
RUN . $LOCALENV ; go mod download

COPY . .
RUN . $LOCALENV ; go build -o /postman-mockserver


FROM alpine:3.12

COPY /docker-entrypoint.sh /docker-entrypoint.sh
RUN ["chmod", "+x", "/docker-entrypoint.sh"]

COPY ./config.yaml /app/config/config.yaml

COPY --from=builder /postman-mockserver /app/postman-mockserver

ENTRYPOINT ["./docker-entrypoint.sh"]

