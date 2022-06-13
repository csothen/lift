#!/bin/bash

DIR=${PWD##*/}

docker rmi $(docker images | grep -e $DIR -e $1 | awk '{ print $3 }')