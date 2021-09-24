#! /bin/sh

echo "#### Starting stream with gpio"
ssh root@192.168.1.11 '/opt/aes/stream-with-gpio /opt/aes/double-spike.aiff'
