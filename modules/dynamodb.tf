# Account table
# Tracks the status of AWS Accounts in our pool
locals {
  // Suffix table names with var.namesapce,
  // unless we're on prod (then no suffix)
  table_suffix = var.namespace == "prod" ? "" : title(var.namespace)
}

resource "aws_dynamodb_table" "accounts" {
  name             = "Accounts${local.table_suffix}"
  read_capacity    = 5
  write_capacity   = 5
  hash_key         = "Id"
  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  global_secondary_index {
    name            = "AccountStatus"
    hash_key        = "AccountStatus"
    projection_type = "ALL"
    read_capacity   = 5
    write_capacity  = 5
  }

  server_side_encryption {
    enabled = true
  }

  # AWS Account ID
  attribute {
    name = "Id"
    type = "S"
  }

  # Status of the Account
  # May be one of:
  #   - LEASED
  #   - READY
  #   - NOT_READY
  attribute {
    name = "AccountStatus"
    type = "S"
  }

  tags = var.global_tags
  /*
  Other attributes:
  - LastModifiedOn (Integer, epoch timestamps)
  - CreatedOn (Integer, epoch timestamps)
  */
}

resource "aws_sns_topic" "accounts" {
  name = "accounts-${var.namespace}"
  tags = var.global_tags
}

module "publish_account_events_lambda" {
  source          = "./lambda"
  name            = "publish_account_events-${var.namespace}"
  namespace       = var.namespace
  description     = "Updated AccountPoolMetrics DB in response to Accounts DB changes"
  global_tags     = var.global_tags
  handler         = "publish_account_events"
  alarm_topic_arn = aws_sns_topic.alarms_topic.arn

  environment = {
    AWS_CURRENT_REGION       = var.aws_region
    ACCOUNT_DB               = aws_dynamodb_table.accounts.id
    ACCOUNT_TOPIC_ARN        = aws_sns_topic.accounts.arn
  }
}

resource "aws_lambda_event_source_mapping" "publish_account_events_from_dynamo_db" {
  event_source_arn  = aws_dynamodb_table.accounts.stream_arn
  function_name     = module.publish_account_events_lambda.name
  batch_size        = 1
  starting_position = "LATEST"
}


resource "aws_dynamodb_table" "leases" {
  name             = "Leases${local.table_suffix}"
  read_capacity    = 5
  write_capacity   = 5
  hash_key         = "AccountId"
  range_key        = "PrincipalId"
  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  server_side_encryption {
    enabled = true
  }

  global_secondary_index {
    name            = "PrincipalId"
    hash_key        = "PrincipalId"
    projection_type = "ALL"
    read_capacity   = 5
    write_capacity  = 5
  }

  global_secondary_index {
    name            = "LeaseStatus"
    hash_key        = "LeaseStatus"
    projection_type = "ALL"
    read_capacity   = 5
    write_capacity  = 5
  }

  global_secondary_index {
    name            = "LeaseId"
    hash_key        = "Id"
    projection_type = "ALL"
    read_capacity   = 5
    write_capacity  = 5
  }

  # AWS Account ID
  attribute {
    name = "AccountId"
    type = "S"
  }

  # Lease status.
  # May be one of:
  # - ACTIVE
  # - INACTIVE
  attribute {
    name = "LeaseStatus"
    type = "S"
  }

  # Principal ID
  attribute {
    name = "PrincipalId"
    type = "S"
  }

  # Lease ID
  attribute {
    name = "Id"
    type = "S"
  }

  tags = var.global_tags
  /*
  Other attributes:
    - LeaseStatusReason (string)
    - CreatedOn (Integer, epoch timestamps)
    - LastModifiedOn (Integer, epoch timestamps)
    - LeaseStatusModifiedOn (Integer, epoch timestamps)
  */
}

resource "aws_dynamodb_table" "usage" {
  name             = "Usage${local.table_suffix}"
  read_capacity    = 5
  write_capacity   = 5
  hash_key         = "StartDate"
  range_key        = "PrincipalId"
  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  server_side_encryption {
    enabled = true
  }

  # User Principal ID
  attribute {
    name = "PrincipalId"
    type = "S"
  }

  # AWS usage cost amount for start date as epoch timestamp
  attribute {
    name = "StartDate"
    type = "N"
  }

  # TTL enabled attribute
  ttl {
    attribute_name = "TimeToLive"
    enabled        = true
  }

  tags = var.global_tags
}
