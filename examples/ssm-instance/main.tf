provider "aws" {
  region = var.region
}

data "aws_availability_zones" "available" {
  state = "available"
}

# Minimal Ephemeral VPC
resource "aws_vpc" "this" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  tags = { Name = "test-ssm-${var.run_id}" }
}

resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.this.id
}

resource "aws_subnet" "this" {
  vpc_id                  = aws_vpc.this.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = data.aws_availability_zones.available.names[0]
}

resource "aws_route_table" "this" {
  vpc_id = aws_vpc.this.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.this.id
  }
}

resource "aws_route_table_association" "this" {
  subnet_id      = aws_subnet.this.id
  route_table_id = aws_route_table.this.id
}

module "ssm_profile" {
  source      = "../../modules/instance-ssm-profile"
  name_prefix = "test-ssm-${var.run_id}-"
}

data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-x86_64"]
  }
}

resource "aws_security_group" "no_ingress" {
  name        = "test-ssm-sg-${var.run_id}"
  description = "No ingress ports allowed"
  vpc_id      = aws_vpc.this.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "this" {
  ami                         = data.aws_ami.amazon_linux.id
  instance_type               = "t3.micro"
  subnet_id                   = aws_subnet.this.id
  iam_instance_profile        = module.ssm_profile.profile_name
  vpc_security_group_ids      = [aws_security_group.no_ingress.id]
  associate_public_ip_address = true # Required for SSM agent to reach API

  tags = {
    Name = "test-ssm-${var.run_id}"
  }
}