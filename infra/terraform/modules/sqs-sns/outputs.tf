output "queue_arns" {
  description = "Map of queue name to ARN (standard + FIFO)"
  value = merge(
    { for k, v in aws_sqs_queue.standard : k => v.arn },
    { for k, v in aws_sqs_queue.fifo : k => v.arn },
  )
}

output "queue_urls" {
  description = "Map of queue name to URL (standard + FIFO)"
  value = merge(
    { for k, v in aws_sqs_queue.standard : k => v.id },
    { for k, v in aws_sqs_queue.fifo : k => v.id },
  )
}

output "dlq_arns" {
  description = "Map of DLQ name to ARN (standard + FIFO)"
  value = merge(
    { for k, v in aws_sqs_queue.standard_dlq : k => v.arn },
    { for k, v in aws_sqs_queue.fifo_dlq : k => v.arn },
  )
}

output "topic_arns" {
  description = "Map of topic name to ARN"
  value = {
    for k, v in aws_sns_topic.topics : k => v.arn
  }
}
