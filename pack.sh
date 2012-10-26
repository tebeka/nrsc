#!/bin/bash
# Pack assets as zip payload in go executable

# Idea from Carlos Castillo (http://bit.ly/SmYXXm)

case "$1" in
    -h | --help ) echo "usage: $(basename $0) EXECTABLE RESOURCE_DIR"; exit;;
esac

if [ $# -ne 2 ]; then
    $0 -h
    exit 1
fi

exe=$1
root=$2

if [ ! -f "${exe}" ]; then
    echo "error: can't find $exe"
    exit 1
fi

if [ ! -d "${root}" ]; then
    echo "error: ${root} is not a directory"
    exit 1
fi

# Exit on 1'st error
set -e

tmp="/tmp/nrsc-$(date +%s).zip"
trap "rm -f ${tmp}" EXIT

(cd "${root}" && zip -r "${tmp}" .)

cat "${tmp}" >> "${exe}"
zip -q -A "${exe}"
