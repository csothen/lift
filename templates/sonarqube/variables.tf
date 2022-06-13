variable "private_key" {
  type=string
  default = "/lift/.ssh/id_rsa"
}

variable "public_key" {
  type=string
  default = "/lift/.ssh/id_rsa.pub"
}