#!/bin/bash

createuser sonar;
echo "ALTER ROLE sonar WITH PASSWORD '$1'; CREATE DATABASE sonar OWNER sonar; \q" | psql;