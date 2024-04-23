
module "security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 5.0"

  name        = local.name
  description = "Complete PostgreSQL example security group"
  vpc_id      = module.vpc.vpc_id

  # ingress
  ingress_with_cidr_blocks = [
    {
      from_port   = 5432
      to_port     = 5432
      protocol    = "tcp"
      description = "PostgreSQL access from within VPC"
      cidr_blocks = module.vpc.vpc_cidr_block
    },
  ]

}

resource "random_password" "rds_password" {
  length = 32
}

module "db" {
  source = "terraform-aws-modules/rds/aws"

  identifier = local.name

  db_subnet_group_name   = module.vpc.database_subnet_group
  vpc_security_group_ids = [module.security_group.security_group_id]

  engine            = "postgres"
  family            = "postgres16"
  engine_version    = "16.1"
  instance_class    = "db.t3.micro"
  allocated_storage = 10

  multi_az                    = false
  publicly_accessible         = false
  create_db_option_group      = false
  create_db_parameter_group   = false
  storage_encrypted           = false
  manage_master_user_password = false


  db_name  = "simple_bank"
  username = "aseerkt"
  password = random_password.rds_password.bcrypt_hash
  port     = "5432"

  deletion_protection = false

  tags = {
    Name = "${local.name}/RDS"
  }
}
