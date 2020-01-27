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
    ACCOUNT_TOPIC_ARN        = aws_sns_topic.lease_locked.arn
  }
}

resource "aws_lambda_event_source_mapping" "publish_account_events_from_dynamo_db" {
  event_source_arn  = aws_dynamodb_table.accounts.stream_arn
  function_name     = module.publish_account_events_lambda.name
  batch_size        = 1
  starting_position = "LATEST"
}

resource "aws_iam_role_policy" "publish_account_events_lambda_dynamo_db" {
  role   = module.publish_account_events_lambda.execution_role_name
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
        "Effect": "Allow",
        "Action": [
            "dynamodb:DescribeStream",
            "dynamodb:GetRecords",
            "dynamodb:GetShardIterator",
            "dynamodb:ListStreams"
        ],
        "Resource": "${aws_dynamodb_table.accounts.stream_arn}"
    },
    {
        "Effect": "Allow",
        "Action": [
            "sns:Publish"
        ],
        "Resource": [
            "*"
        ]
    },
    {
        "Effect": "Allow",
        "Action": [
            "sqs:Publish"
        ],
        "Resource": [
            "*"
        ]
    }
  ]
}
POLICY
}
