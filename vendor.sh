#!/bin/bash

display_help()
{
    echo
    echo "usage: $0 -u url"
    echo "option:"
    echo "-u url(github)"
    echo
    echo
    exit 0
}

CUR_PATH=$(cd `dirname $0`; pwd)
# parse options usage: $0 -c config.ini -m flag -b base_path -r remote_path -u user_name -p pwd -t
while getopts 'u:' OPT; do
    case $OPT in
        u)
            git_url="$OPTARG";;
        ?)
            display_help
    esac
done

if [ ! -n "$git_url" ]; then
    echo "url must be set up"
    exit 0
fi

exist=`which govendor | wc -l`
if [ $exist == 0 ]; then
    echo "please install govendor"
    go get -u github.com/kardianos/govendor
fi
go get $git_url
govendor list -v list
govendor add +external
govendor list -v list
