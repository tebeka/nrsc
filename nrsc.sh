#!/bin/bash
# Pack assets as zip payload in go executable

# Idea from Carlos Castillo (http://bit.ly/SmYXXm)

case "$1" in
    -h | --help )
        echo "usage: $(basename $0) EXECTABLE RESOURCE_DIR [ZIP OPTIONS]";
        exit;;
    --version )
        echo "nrsc version 0.3.1"; exit;;
esac

if [ $# -lt 2 ]; then
    $0 -h
    exit 1
fi

exe=$1
shift
root=$1
shift

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

# Create zip file
(cd "${root}" && zip -r "${tmp}" . $@)

# Append zip to executable
cat "${tmp}" >> "${exe}"
# Fix zip offset in file
zip -q -A "${exe}"
