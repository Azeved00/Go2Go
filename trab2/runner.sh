#!/usr/bin/env bash


work_dir="$PWD/trab2"
go="GOWORK=\"$work_dir/go.work\" go"

peer0="$go run $work_dir/peer/main.go -p 8081 -n localhost:8083 localhost"
peer1="$go run $work_dir/peer/main.go -p 8082 -n localhost:8083 localhost"
peer2="$go run $work_dir/peer/main.go -p 8083 -n localhost:8084 localhost"
peer3="$go run $work_dir/peer/main.go -p 8084 -n localhost:8085 localhost"
peer5="$go run $work_dir/peer/main.go -p 8085 localhost"
peer4="$go run $work_dir/peer/main.go -p 8086 -n localhost:8084 localhost"

# Check the number of panes in the current window
pane_count=$(tmux list-panes | wc -l)

#create panes if needed
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
tmux send-keys -t 0 "$peer0" Enter
tmux send-keys -t 1 "$peer1" Enter
tmux send-keys -t 2 "$peer2" Enter
tmux send-keys -t 3 "$peer3" Enter
tmux send-keys -t 4 "$peer4" Enter
tmux send-keys -t 5 "$peer5" Enter
