resource "aws_ecr_repository" "this" {
  name = local.name

  image_tag_mutability = "IMMUTABLE"
  force_delete         = true

  tags = {
    Name = "${local.name}/ECR"
  }
}
