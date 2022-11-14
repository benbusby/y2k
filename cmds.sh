#!/bin/sh

function set_file_time() {
    case "$(uname -s)" in
        # MacOS uses a different command for modifying file time
        Darwin*) touch -d $(date -r $1 '+%Y-%m-%dT%H:%M:%S') $2;;
        *) touch --date=@$1 $2
    esac
}
