resource "aws_db_subnet_group" "database" {
  depends_on =[aws_subnet.dbsubnet]
  name       = "aws_db_subnet_group-demo-${random_string.role.id}"
  subnet_ids = aws_subnet.dbsubnet[*].id
  tags = {
    Name = "DB subnet group"
  }
}

resource "aws_db_instance" "default" {
  depends_on        = [aws_security_group.dbsg]
  storage_type      = "gp2"
  #DB Network
  vpc_security_group_ids = [aws_security_group.dbsg.id]
  db_subnet_group_name  =  aws_db_subnet_group.database.name
  publicly_accessible = false
  # Storage Allocation
  max_allocated_storage = 20
  allocated_storage    = 10
  #Type, version of DB and Instance Class to use
  engine               = "postgres"
  engine_version       = "11"
  instance_class       = "db.t3.micro"
  #Credentials
  name                 = "postgres"
  username             = "postgres"
  password             = "postgres"
  #Setting this true so that there will be no problem while destroying the Infrastructure as it won't create snapshot
  skip_final_snapshot  = true
  auto_minor_version_upgrade = true
}

resource "aws_security_group" "dbsg" {
  depends_on =[aws_subnet.dbsubnet]
  name        = "db"
  description = "security group for db"
  vpc_id      = aws_vpc.demo.id


  # Allowing traffic only for Postgres and that too from same VPC only.
  ingress {
    description = "POSTGRES"
    from_port   = 5432
    to_port     = 5432
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
    Name = "db-sg"
  }
}
