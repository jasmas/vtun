#!bin/bash

sudo killall vtun
sudo ../bin/vtun -S -l=:3001 -c=172.16.0.1/24 &
echo "started!"
