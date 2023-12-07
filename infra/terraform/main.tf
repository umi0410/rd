terraform {
  backend "s3" {
    bucket = "rd-infra-cmqtaw5mcmek"
    key    = "terraform/rd.json"
    region = "ap-northeast-2"
  }
}

resource "aws_dynamodb_table" "rd" {
  name         = "rd-${terraform.workspace}"
  billing_mode = "PAY_PER_REQUEST"
  # partition key
  hash_key = "PK"
  # sort key
  range_key = "SK"
  attribute {
    name = "PK"
    type = "S"
  }
  attribute {
    name = "SK"
    type = "S"
  }

  tags = {
    Name = "rd-${terraform.workspace}"
  }
}
