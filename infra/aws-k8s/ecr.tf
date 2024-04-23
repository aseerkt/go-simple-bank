resource "aws_ecr_repository" "this" {
  name         = local.name
  force_delete = true
  tags = {
    Name = "${local.name}/ECR"
  }
}
