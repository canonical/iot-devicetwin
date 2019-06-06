FROM golang:1.12 as builder1
COPY . ./src/github.com/CanonicalLtd/iot-devicetwin
WORKDIR /go/src/github.com/CanonicalLtd/iot-devicetwin
RUN ./get-deps.sh
RUN CGO_ENABLED=1 GOOS=linux go build -a -o /go/bin/devicetwin -ldflags='-extldflags "-static"' cmd/devicetwin/main.go

# Copy the built applications to the docker image
FROM ubuntu:18.04
WORKDIR /srv
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY --from=builder1 /go/bin/devicetwin /srv/devicetwin
RUN mkdir -p /srv/certs

# TODO: need the certs from k8s config
COPY --from=builder1 /go/src/github.com/CanonicalLtd/iot-devicetwin/certs/* /srv/certs/

# Set params from the environment variables
ARG DRIVER="postgres"
ARG DATASOURCE="dbname=devicetwin sslmode=disable"
ARG PORT="8040"
ARG MQTTURL="localhost"
ARG MQTTPORT="8883"
ARG CONFIGDIR="/srv/certs"
ENV DRIVER="${DRIVER}"
ENV DATASOURCE="${DATASOURCE}"
ENV PORT="${PORT}"
ENV MQTTURL="${MQTTURL}"
ENV MQTTPORT="${MQTTPORT}"
ENV CONFIGDIR="${CONFIGDIR}"

EXPOSE 8040
ENTRYPOINT /srv/devicetwin -port $PORT -driver $DRIVER -datasource "${DATASOURCE}" -mqtturl $MQTTURL -mqttport $MQTTPORT -configdir $CONFIGDIR