#!bin/bash
SERVER="i.qtun.me:3001"
sudo killall vtun
sudo ./bin/vtun -s=$SERVER -c=172.16.0.2/24 &
echo "started!"
