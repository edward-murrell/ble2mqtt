# Build on the host's platform and cross-compile, so multi-arch builds need no emulation.
FROM --platform=$BUILDPLATFORM golang:1.21 AS build
ARG TARGETOS TARGETARCH TARGETVARIANT
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOARM="${TARGETVARIANT#v}" CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -o /out/ble2mqtt .

# Runs as root so BlueZ accepts scan requests over the host's D-Bus socket.
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/ble2mqtt /usr/sbin/ble2mqtt
ENV BLE2MQTT_CONFIG=/etc/ble2mqtt/ble2mqtt.conf
# Run with: -v /run/dbus:/run/dbus:ro -v ./ble2mqtt.conf:/etc/ble2mqtt/ble2mqtt.conf:ro
ENTRYPOINT ["/usr/sbin/ble2mqtt"]
