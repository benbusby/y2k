#!/bin/sh

source "$(dirname $(readlink -f "$0"))/../../cmds.sh"

set_file_time 401020310 01.y2k
set_file_time 10030100  02.y2k
