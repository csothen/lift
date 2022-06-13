#!/bin/sh

/lift/main server "$@" &
/lift/main observer "$@" &
wait