#! /bin/bash

wget https://gitlab.com/MMw_Unibo/platformeng/federatedlearning-processes/-/archive/main/federatedlearning-processes-main.zip

yes | unzip federatedlearning-processes-main.zip

rm federatedlearning-processes-main.zip


cd ./federatedlearning-processes-main

python3 -m venv .venv

source .venv/bin/activate
pip3 install -r requirements.txt
python3 mqttsub.py --port 3001 > /dev/null 2>&1 &

echo "Federated Learning Process started"
exit 0
