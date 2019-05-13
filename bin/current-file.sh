#!/usr/bin/env bash

# lsof +D /Volumes/Leeroy/Danas_Photos/Dashcam_Scratchpad/Dashcam/_all/ | grep MP4 | sed -e 's/^.*2018_/2018_/'
lsof -p $(ps aux | grep -i OBS | grep -v grep | awk '{print $2}') | grep MP4 | sed -e 's/^.*2018_/2018_/'
