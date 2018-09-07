#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

porcelain_status="#($CURRENT_DIR/porcelain -tmux -path '#{pane_current_path}')"
porcelain_interpolation_string="\#{porcelain}"

do_interpolation() {
    local string="$1"
    local interpolated="${string/$porcelain_interpolation_string/$porcelain_status}"
    echo "$interpolated"
}

update_tmux_option() {
    local option="$1"
    local option_value="$(tmux show-option -gqv "$option")"
    local new_option_value="$(do_interpolation "$option_value")"
    tmux set-option -gq "$option" "$new_option_value"
}

main() {
    if ! [[ -f "$CURRENT_DIR/porcelain" ]]
    then
        curl -sL "$(curl -s https://api.github.com/repos/robertgzr/porcelain/releases/latest | grep "browser_download_url" | grep "$(uname | tr '[:upper:]' '[:lower:]')"| cut -d\" -f4)" |\
            tar -C "$CURRENT_DIR" -xzf - porcelain &>/dev/null
    fi
    # try again
    if [[ -f "$CURRENT_DIR/porcelain" ]]
    then
        update_tmux_option "status-right"
        update_tmux_option "status-left"
    fi
}
main
