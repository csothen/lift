#!/bin/bash

sudo apt-add-repository ppa:ansible/ansible -y
sudo apt update
sudo apt install ansible -y

sleep 15s 
ansible-playbook /tmp/ansible/ansible.yaml -i /tmp/ansible/inventory