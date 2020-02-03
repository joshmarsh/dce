/*
Keep a table of basic account pool metrics up to date using events from the accounts queue
*/
resource "aws_dynamodb_table" "account_pool_metrics" {
  name             = "AccountPoolMetrics${local.table_suffix}"
  read_capacity    = 5
  write_capacity   = 5
  hash_key         = "Id"
  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  server_side_encryption {
    enabled = true
  }

  attribute {
    name = "Id"
    type = "S"
  }

  tags = var.global_tags
}

resource "aws_sqs_queue" "account_pool_metrics" {
  name = "account_pool_metrics-${var.namespace}"
  visibility_timeout_seconds = 3600
  tags = var.global_tags
}

resource "aws_sns_topic_subscription" "account_pool_metrics_subscription" {
  topic_arn = aws_sns_topic.accounts.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.account_pool_metrics.arn
}

module "update_account_pool_metrics_lambda" {
  source          = "./lambda"
  name            = "update_account_pool_metrics-${var.namespace}"
  namespace       = var.namespace
  description     = "Updated AccountPoolMetrics DB in response to Accounts DB changes"
  global_tags     = var.global_tags
  handler         = "update_account_pool_metrics"
  alarm_topic_arn = aws_sns_topic.alarms_topic.arn

  environment = {
    AWS_CURRENT_REGION       = var.aws_region
    ACCOUNT_DB               = aws_dynamodb_table.account_pool_metrics.id
  }
}

resource "aws_lambda_event_source_mapping" "update_account_pool_metrics_from_accounts_sqs" {
  event_source_arn  = aws_sqs_queue.account_pool_metrics.arn
  function_name     = module.update_account_pool_metrics_lambda.name
}