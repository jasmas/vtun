#!bin/bash
SERVER="127.0.0.1:3001"
sudo killall vtun
sudo ./bin/vtun -s=$SERVER -c=172.16.0.2/24 &
echo "started!"
