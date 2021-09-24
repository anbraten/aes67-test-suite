#! /bin/sh

echo "#### Starting alsaloop"

ssh root@192.168.1.12 'nice -n -10 alsaloop -c 2 -r 48000 -f S16_LE -C plughw:RAVENNA -P plughw:C8CH'
