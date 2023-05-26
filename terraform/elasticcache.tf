resource "aws_elasticache_subnet_group" "redis" {
  depends_on = [aws_subnet.redissubnet]
  name       = "redis-subnet-${random_string.cluster.id}"
  subnet_ids = aws_subnet.redissubnet[*].id
}

resource "aws_elasticache_cluster" "demo" {
  depends_on           = [aws_security_group.redissg]
  cluster_id           = "demo-redis-cluster-${random_string.cluster.id}"
  engine               = "redis"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis5.0"
  engine_version       = "5.0.6"
  port                 = 6379
  apply_immediately = true
  subnet_group_name = aws_elasticache_subnet_group.redis.name
  security_group_ids = [aws_security_group.redissg.id]
}

resource "aws_security_group" "redissg" {
  depends_on = [aws_subnet.redissubnet]
  name        = "redis-sg-${random_string.cluster.id}"
  description = "security group for redis"
  vpc_id      = aws_vpc.demo.id


  # Allowing traffic only for Postgres and that too from same VPC only.
  ingress {
    description = "Redis"
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [var.cluster_ipv4_cidr]
  }


  # Allowing all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "redis-sg"
  }
}
