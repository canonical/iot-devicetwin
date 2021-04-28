data "kubectl_file_documents" "twin_manifests" {
  content = templatefile(
    "${path.module}/pods/k8s-devicetwin.yaml",
    { IMAGE = "${local.docker_namespace}/iot-devicetwin:${local.docker_tag}" }
  )
}

resource "kubectl_manifest" "twin_identity" {
  override_namespace = local.namespace
  count              = length(data.kubectl_file_documents.twin_manifests.documents)
  yaml_body          = element(data.kubectl_file_documents.twin_manifests.documents, count.index)
  wait               = true
}

resource "null_resource" "rollout_twin" {
  depends_on = [
    kubectl_manifest.twin_identity
  ]
  triggers = {
    "timestamp" = timestamp()
  }
  provisioner "local-exec" {
    command = <<EOF
      aws eks update-kubeconfig --region ${local.region} --name ${local.eks_cluster};
      kubectl rollout restart deploy devicetwin -n=${local.namespace};
    EOF
  }
}