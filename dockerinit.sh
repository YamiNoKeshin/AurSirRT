#!/bin/sh

./main -zconnection=$(ifconfig eth0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}'):5555
