#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PORCELAIN_BIN="$CURRENT_DIR/porcelain"

porcelain_status="#($PORCELAIN_BIN -tmux -path '#{pane_current_path}')"
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

install() {
    update_tmux_option "status-right"
    update_tmux_option "status-left"
}

main() {
    if ! [[ -f "$PORCELAIN_BIN" ]]
    then
        curl -sL "$(curl -s https://api.github.com/repos/robertgzr/porcelain/releases/latest | grep "browser_download_url" | grep "$(uname | tr '[:upper:]' '[:lower:]')"| cut -d\" -f4)" |\
            tar -C "$CURRENT_DIR" -xzf - porcelain &>/dev/null
    fi
    # try again
    if [[ -f "$PORCELAIN_BIN" ]]
    then
        install
        exit 0
    fi
}
main
