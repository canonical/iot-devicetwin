#!/bin/sh


CA_ORG="/O=IoT Management/emailAddress=${EMAIL:-admin@example.com}"
CA_DN="/CN=mqtt${CA_ORG}"
MQTT_DN="/CN=mqtt$CA_ORG"
TWIN_DN="/CN=mqtt$CA_ORG"

# Certificate authority
openssl req -newkey rsa:2048 -x509 -nodes -sha512 -days 3650 -extensions v3_ca -keyout ca.key -out ca.crt -subj "${CA_DN}"


# MQTT broker certificate
openssl genrsa -out mqtt.key 2048
openssl req -new -sha512 -out mqtt.csr -key mqtt.key -subj "${MQTT_DN}"

cat > "./mqtt.cnf" << EOF
[v3_req]
basicConstraints        = critical,CA:FALSE
nsCertType              = server
subjectKeyIdentifier    = hash
authorityKeyIdentifier  = keyid,issuer:always
keyUsage                = nonRepudiation, digitalSignature, keyEncipherment
EOF
openssl x509 -req -sha512 -in mqtt.csr -CA ca.crt -CAkey ca.key -CAcreateserial -CAserial ca.srl -out mqtt.crt -days 365 -extfile "./mqtt.cnf" -extensions v3_req


# Device twin certificate
openssl genrsa -out devicetwin.key 2048
openssl req -new -sha512 -out devicetwin.csr -key devicetwin.key -subj "${TWIN_DN}"

cat > "./devicetwin.cnf" << EOF
[v3_req]
basicConstraints        = critical,CA:FALSE
nsCertType              = client
extendedKeyUsage        = clientAuth
subjectKeyIdentifier    = hash
authorityKeyIdentifier  = keyid,issuer:always
keyUsage                = digitalSignature, keyEncipherment, keyAgreement
EOF
openssl x509 -req -sha512 -in devicetwin.csr -CA ca.crt -CAkey ca.key -CAcreateserial -CAserial ca.srl -out devicetwin.crt -days 365 -extfile "./devicetwin.cnf" -extensions v3_req

# Generate the kubernetes secrets
cat > "./devicetwin.yaml" << EOF
apiVersion: v1
kind: Secret
metadata:
  name: devicetwin-certs
data:
  # base64 encoded X509 certificate files
  ca.crt:     "$(base64 -w0 ca.crt)"
  server.crt: "$(base64 -w0 devicetwin.crt)"
  server.key: "$(base64 -w0 devicetwin.key)"
---
EOF

cat > "./mqtt.yaml" << EOF
apiVersion: v1
kind: Secret
metadata:
  name: mqtt-certs
data:
  # base64 encoded X509 certificate files
  ca.crt:     "$(base64 -w0 ca.crt)"
  server.crt: "$(base64 -w0 mqtt.crt)"
  server.key: "$(base64 -w0 mqtt.key)"
---
EOF

cat > "./identity.yaml" << EOF
apiVersion: v1
kind: Secret
metadata:
  name: identity-certs
data:
  # base64 encoded X509 certificate files
  ca.crt: "$(base64 -w0 ca.crt)"
  ca.key: "$(base64 -w0 ca.key)"
---
EOF

# Clean up
rm -- *.cnf *.csr

echo "Use the *.yaml files for deploying the X509 certificates as kubernetes secrets"
