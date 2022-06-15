terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "3.74.0"
    }
  }
}

provider "aws" {
  profile = "default"
  region = "eu-west-1"
}

resource "aws_key_pair" "($ name $)" {
  key_name = "($ name $)-key-pair"
  public_key = file("${var.public_key}")
}

resource "aws_security_group" "($ name $)" {
    name        = "($ name $)-sg"

    ingress {
      description = "Access Jenkins"
      from_port   = 8080
      to_port     = 8080
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
      description = "SSH from everywhere"
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }

    tags = {
        Name = "($ name $)-sg"
    }
}

resource "aws_instance" "($ name $)" {
  count = "($ count $)"
  // retrieved from https://cloud-images.ubuntu.com/locator/ec2/
  // and it refers to an eu-west-1 Ubuntu 18.04LTS amd64 machine
  ami           = "ami-0ce48dd7b483b8402"
  instance_type = "t2.micro"
  key_name = aws_key_pair.($ name $).key_name
  
  security_groups = [aws_security_group.($ name $).name]

  tags = {
    "Name" = "($ name $)-${count.index}"
  }
}

resource "null_resource" "run_ansible" {
  count = "($ count $)"

  connection {
    type = "ssh"
    user = "ubuntu"
    host = aws_instance.($ name $)[count.index].public_ip
    private_key="${file("${var.private_key}")}"
    agent = false
    timeout = "3m"
  }

  provisioner "file" {
    source = "./files"
    destination = "/tmp"
  }

  provisioner "file" {
    source = "./ansible"
    destination = "/tmp"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/files/run-ansible.sh",
      "/tmp/files/run-ansible.sh"
    ]
  }
}

output "public_ips" {
  description = "Public IPs of the instances"
  value = "${aws_instance.($ name $).*.public_ip}"
}