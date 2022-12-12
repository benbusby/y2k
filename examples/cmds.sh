#!/bin/sh

function set_file_time() {
    case "$(uname -s)" in
        # MacOS uses a different command for modifying file time
        Darwin*) gtouch --date=@$1 $2;;
        *) touch --date=@$1 $2
    esac
}
