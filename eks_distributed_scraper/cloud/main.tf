resource "aws_iam_role" "iam-role-eks-scraper-cluster" {
  name               = "eks-scraper-cluster"
  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
   "Effect": "Allow",
   "Principal": {
    "Service": "eks.amazonaws.com"
   },
   "Action": "sts:AssumeRole"
   }
  ]
 }
EOF
}

resource "aws_iam_role_policy_attachment" "eks-scraper-cluster-AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.iam-role-eks-scraper-cluster.name
}

resource "aws_iam_role_policy_attachment" "eks-scraper-cluster-AmazonEKSServicePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
  role       = aws_iam_role.iam-role-eks-scraper-cluster.name
}

resource "aws_vpc" "scraper-vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = "true"
  enable_dns_hostnames = "false"
  enable_classiclink   = "false"
  instance_tenancy     = "default"

  tags = {
    Name = "testVpcName"
  }
}

resource "aws_subnet" "scraper-pub-subnet-1" {
  vpc_id                  = aws_vpc.scraper-vpc.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = "true"
  availability_zone       = var.AWS_DEFAULT_ZONE

  tags = {
    Name = "public-eks-subnet-1"
  }
}

resource "aws_subnet" "scraper-priv-subnet-1" {
  vpc_id            = aws_vpc.scraper-vpc.id
  cidr_block        = "10.0.129.0/24"
  availability_zone = var.AWS_DEFAULT_ZONE

  tags = {
    Name = "private-eks-subnet-1"
  }
}

resource "aws_subnet" "scraper-priv-subnet-2" {
  vpc_id            = aws_vpc.scraper-vpc.id
  cidr_block        = "10.0.130.0/24"
  availability_zone = var.AWS_SECONDARY_ZONE

  tags = {
    Name = "private-eks-subnet-2"
  }
}

resource "aws_internet_gateway" "scraper-igw" {
  vpc_id = aws_vpc.scraper-vpc.id
  tags = {
    Name = "scraper-igw"
  }
}

resource "aws_route_table" "scraper-public-crt" {
  vpc_id = aws_vpc.scraper-vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.scraper-igw.id
  }
  tags = {
    Name = "scraper-publict-crt"
  }
}

resource "aws_route_table_association" "scraper-crta-pub-subnet-1" {
  subnet_id      = aws_subnet.scraper-pub-subnet-1.id
  route_table_id = aws_route_table.scraper-public-crt.id
}

resource "aws_security_group" "eks-scraper-SG" {
  name   = "SG-eks-scraper-cluster"
  vpc_id = aws_vpc.scraper-vpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # ingress {
  #   from_port   = 22
  #   to_port     = 22
  #   protocol    = "tcp"
  #   cidr_blocks = ["0.0.0.0/0"]
  # }

  # ingress {
  #   from_port   = 80
  #   to_port     = 80
  #   protocol    = "tcp"
  #   cidr_blocks = ["0.0.0.0/0"]
  # }
}

resource "aws_eks_cluster" "eks_scraper_cluster" {
  name     = "eks_scraper_cluster"
  role_arn = aws_iam_role.iam-role-eks-scraper-cluster.arn
  version  = 1.21

  vpc_config {
    security_group_ids = ["${aws_security_group.eks-scraper-SG.id}"]
    subnet_ids         = ["${aws_subnet.scraper-priv-subnet-1.id}", "${aws_subnet.scraper-priv-subnet-2.id}"]
  }

  depends_on = [
    aws_iam_role_policy_attachment.eks-scraper-cluster-AmazonEKSClusterPolicy,
    aws_iam_role_policy_attachment.eks-scraper-cluster-AmazonEKSServicePolicy,
  ]
}

resource "aws_iam_role" "scraper-node-role" {
  name = "eks-node-group"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "ec2.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "AmazonEKSWorkerNodePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.scraper-node-role.name
}

resource "aws_iam_role_policy_attachment" "AmazonEKS_CNI_Policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.scraper-node-role.name
}

resource "aws_iam_role_policy_attachment" "AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.scraper-node-role.name
}

resource "aws_eks_node_group" "node" {
  cluster_name    = aws_eks_cluster.eks_scraper_cluster.name
  node_group_name = "scraper_node_group_1"
  node_role_arn   = aws_iam_role.scraper-node-role.arn
  subnet_ids      = ["${aws_subnet.scraper-priv-subnet-1.id}", "${aws_subnet.scraper-priv-subnet-2.id}"]
  instance_types  = ["t3.micro", "t2.micro"]

  scaling_config {
    desired_size = 1
    max_size     = 2
    min_size     = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.AmazonEC2ContainerRegistryReadOnly,
  ]
}

output "tags" {
  value = var.DEFAULT_TAGS
}