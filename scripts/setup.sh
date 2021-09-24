#! /bin/sh

ssh root@192.168.1.10 'killall nodejs'

ssh root@192.168.1.11 'killall nodejs'
ssh root@192.168.1.11 '/etc/init.d/aes start'

ssh root@192.168.1.12 'killall nodejs'
ssh root@192.168.1.12 '/etc/init.d/aes start'
