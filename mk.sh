#! /usr/bin/env bash

COLOR="\e[96m\e[1m"
ENDCOLOR="\e[0m"

echo -e "${COLOR}BUILDING marlinctl2 with version $1 ${ENDCOLOR}"
export MARLINCTL2BUILDVERSIONSTRING=$1
make release

echo -e "${COLOR}COPYING marlinctl2 to /usr/local/bin/ ${ENDCOLOR}"
make install
