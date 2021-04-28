data "kubectl_file_documents" "mosquitto_manifests" {
  content = templatefile(
    "${path.module}/pods/k8s-mosquitto.yaml",
    { IP_WHITELIST = local.ip_whitelist_public }
  )
}

resource "kubectl_manifest" "iot_mosquitto" {
  override_namespace = local.namespace
  count              = length(data.kubectl_file_documents.mosquitto_manifests.documents)
  yaml_body          = element(data.kubectl_file_documents.mosquitto_manifests.documents, count.index)
  wait               = true
}

data "kubernetes_service" "iot_mosquitto" {
  depends_on = [
    kubectl_manifest.iot_mosquitto
  ]
  metadata {
    name      = "mqtt-lb"
    namespace = local.namespace
  }
}

data "aws_lb" "iot_mosquitto" {
  depends_on = [
    kubectl_manifest.iot_mosquitto
  ]
  name = split("-", data.kubernetes_service.iot_mosquitto.status.0.load_balancer.0.ingress.0.hostname)[0]
}

resource "aws_route53_record" "iot_mosquitto" {
  zone_id = local.hosted_zone
  name    = local.mqtt_domain
  type    = "A"
  alias {
    name                   = data.aws_lb.iot_mosquitto.dns_name
    zone_id                = data.aws_lb.iot_mosquitto.zone_id
    evaluate_target_health = true
  }
}

resource "null_resource" "rollout_mqtt" {
  depends_on = [
    aws_route53_record.iot_mosquitto
  ]
  triggers = {
    "timestamp" = timestamp()
  }
  provisioner "local-exec" {
    command = <<EOF
      aws eks update-kubeconfig --region ${local.region} --name ${local.eks_cluster};
      kubectl rollout restart deploy mqtt -n=${local.namespace};
    EOF
  }
}