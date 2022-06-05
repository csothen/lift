#!/bin/bash

DIR=$(dirname "${BASH_SOURCE[0]}")
INIT=$(cat $DIR/init.template.sql)

echo "$INIT" | sed "s/(db_user)/$1/;s/(db_name)/$2/" > $DIR/init.sql