terraform {
  backend "s3" {
    bucket = "rd-infra-cmqtaw5mcmek"
    key    = "terraform/rd.json"
    region = "ap-northeast-2"
  }
}

resource "aws_dynamodb_table" "alias" {
  name         = "rd-aliases-${terraform.workspace}"
  billing_mode = "PAY_PER_REQUEST"
  # partition key
  hash_key = "group"
  attribute {
    name = "group_alias"
    type = "S"
  }

  tags = {
    Name = "rd-aliases-${terraform.workspace}"
  }
}

resource "aws_dynamodb_table" "group" {
  name         = "rd-groups-${terraform.workspace}"
  billing_mode = "PAY_PER_REQUEST"
  # partition key
  hash_key = "group"

  attribute {
    name = "group"
    type = "S"
  }

  tags = {
    Name = "rd-groups-${terraform.workspace}"
  }
}

resource "aws_dynamodb_table" "user" {
  name         = "rd-users-${terraform.workspace}"
  billing_mode = "PAY_PER_REQUEST"
  # partition key
  hash_key = "username"

  attribute {
    name = "username"
    type = "S"
  }

  tags = {
    Name = "rd-users-${terraform.workspace}"
  }
}

resource "aws_dynamodb_table" "alias_hit_event" {
  name         = "rd-alias-hit-events-${terraform.workspace}"
  billing_mode = "PAY_PER_REQUEST"
  # partition key
  hash_key  = "group_alias"
  range_key = "created_at"

  attribute {
    name = "group_alias"
    type = "S"
  }

  attribute {
    name = "created_at"
    type = "S"
  }

  tags = {
    Name = "rd-alias-hit-events-${terraform.workspace}"
  }
}
