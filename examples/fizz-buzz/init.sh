#!/bin/sh

source "$(dirname $(readlink -f "$0"))/../cmds.sh"

set_file_time 502080901.040609262 01.y2k
set_file_time 960808010.402212626 02.y2k
set_file_time 905000187.919775188 03.y2k
set_file_time 912106121.310071111 04.y2k
set_file_time 961402159.274200061 05.y2k
set_file_time 940139294.200061401 06.y2k
set_file_time 959284200.009210000 07.y2k

# Shorter version with "f" and "b" instead of
# "fizz" and "buzz"
#set_file_time 812106121.310071111 01.y2k
#set_file_time 161402159.162004200 02.y2k
#set_file_time 106140139.160042000 03.y2k
#set_file_time 161401591.200420009 04.y2k
#set_file_time 921000000.000000000 05.y2k
