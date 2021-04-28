output "mqtt_url" {
  description = "The url of the MQTT Service"
  value       = aws_route53_record.iot_mosquitto.fqdn
}