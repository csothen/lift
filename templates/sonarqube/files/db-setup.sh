#!/bin/bash

createuser sonar;
echo "ALTER USER sonar WITH ENCRYPTED password '$1'; CREATE DATABASE sonar OWNER sonar; \q" | psql;

admin_user="admin"
admin_salt=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 40 ; echo '')
admin_password="$2"
admin_salted_password=$(echo -n "--${admin_salt}--${admin_password}--" | sha1sum | awk '{print $1}')

echo "UPDATE users SET crypted_password='$admin_salted_password', salt='$admin_salt', hash_method='SHA1' WHERE login='$admin_user'" | psql