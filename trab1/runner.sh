#!/usr/bin/env bash

work_dir="$PWD/trab1"
go="GOWORK=\"$work_dir/go.work\" go"

server="$go run $work_dir/server/main.go 8080"

peer1="$go run $work_dir/peer/main.go -p 8081 -n 8082 localhost"
peer2="$go run $work_dir/peer/main.go -p 8082 -n 8083 localhost"
peer3="$go run $work_dir/peer/main.go -p 8083 -n 8084 localhost"
peer4="$go run $work_dir/peer/main.go -p 8084 -n 8085 localhost"
peer5="$go run $work_dir/injector/main.go localhost:8081 && $go run $work_dir/peer/main.go -p 8085 -n 8081 localhost"

# Check the number of panes in the current window
pane_count=$(tmux list-panes | wc -l)

# If there are fewer than 3 panes, create more
if [ "$pane_count" -lt 6 ]; then
    tmux split-window -v    
    tmux split-window -h   
    tmux split-window -h  
    tmux select-pane -t 0
    tmux split-window -h 
    tmux select-pane -t 2
    tmux split-window -h 
    tmux select-layout tiled
fi


# Send the commands to each of the 3 panes
tmux send-keys -t 0 "$server" Enter
tmux send-keys -t 1 "$peer1" Enter
tmux send-keys -t 2 "$peer2" Enter
tmux send-keys -t 3 "$peer3" Enter
tmux send-keys -t 4 "$peer4" Enter
sleep 2s
tmux send-keys -t 5 "$peer5" Enter
