#!/usr/bin/env bash

# Commands you want to run
cmd1="go run $ROOT/trab1/peer/main.go 8081 8082"
cmd2="go run $ROOT/trab1/server/main.go 8080"
cmd3="go run $ROOT/trab1/injector/main.go 8081 && go run $ROOT/trab1/peer/main.go 8082 8081"

# Check the number of panes in the current window
pane_count=$(tmux list-panes | wc -l)

# If there are fewer than 3 panes, create more
if [ "$pane_count" -lt 3 ]; then
    tmux select-layout even-horizontal  # Set a layout to make room for new panes
    while [ "$pane_count" -lt 3 ]; do
        tmux split-window -h
        pane_count=$((pane_count + 1))
    done
    tmux select-layout tiled  # Rearrange to a 3-pane layout
fi

export $GOWORK=$ROOT/trab1/go.work

# Send the commands to each of the 3 panes
tmux send-keys -t 0 "$cmd1" Enter
tmux send-keys -t 1 "$cmd2" Enter
sleep 2s
tmux send-keys -t 2 "$cmd3" Enter
