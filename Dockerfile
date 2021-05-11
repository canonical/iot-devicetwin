FROM golang:1.16 as builder1
COPY . /iot-devicetwin
WORKDIR /iot-devicetwin
RUN CGO_ENABLED=1 GOOS=linux go build -a -o /go/bin/devicetwin -ldflags='-extldflags "-static"' cmd/devicetwin/main.go

# Copy the built applications to the docker image
FROM ubuntu:18.04
WORKDIR /srv
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY --from=builder1 /go/bin/devicetwin /srv/devicetwin

# Note: /srv/certs is handled by the k8s file

# Set params from the environment variables
ARG DRIVER="postgres"
ARG DATASOURCE="dbname=devicetwin sslmode=disable"
ARG PORT="8040"
ARG MQTTURL="localhost"
ARG MQTTPORT="8883"
ARG CERTSDIR="/srv/certs"
ARG CONFIGDIR="/srv/config"
ENV DRIVER="${DRIVER}"
ENV DATASOURCE="${DATASOURCE}"
ENV PORT="${PORT}"
ENV MQTTURL="${MQTTURL}"
ENV MQTTPORT="${MQTTPORT}"
ENV CERTSDIR="${CERTSDIR}"
ENV CONFIGDIR="${CONFIGDIR}"

EXPOSE 8040
ENTRYPOINT /srv/devicetwin -port $PORT -driver $DRIVER -datasource "${DATASOURCE}" -mqtturl $MQTTURL -mqttport $MQTTPORT -certsdir $CERTSDIR -configdir $CONFIGDIR
