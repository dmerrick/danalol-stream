#!/usr/bin/env bash

# this script finds the currently-playing video file

#TODO: error message if more than 1 result from search

# SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
# lsof -p $(cat ${SCRIPTPATH}/../OBS/OBS.pid) | grep -i '\.MP4' | sed -e 's/^.*2018_/2018_/'

# echo "2018_0522_201919_010_opt.MP4"
# ssh void -- "cd other_projects/danalol-stream && /Users/dmerrick/other_projects/danalol-stream/bin/current-file.sh"
lsof -p $(cat OBS/OBS.pid) | grep -i '\.MP4' | sed -e 's/^.*2018_/2018_/'
