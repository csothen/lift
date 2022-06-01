terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

resource "aws_instance" "($ name $)" {
  // retrieved from https://cloud-images.ubuntu.com/locator/ec2/
  // and it refers to an eu-west-1 Ubuntu 18.04LTS amd64 machine
  ami           = "ami-0ce48dd7b483b8402"
  instance_type = "t2.micro"
  count = "($ count $)"

  tags = {
    "Name" = "($ name $)"
  }

  provisioner "local-exec" {
    command = "ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -u root -i '${self.ipv4_address} ansible.yml"
  }

  output "ec2_global_ips" {
    value = ["${aws_instance.main.*.public_ip}"]
  }
}