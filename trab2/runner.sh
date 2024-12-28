#!/usr/bin/env bash

go="GOWORK=\"$ROOT/trab2/go.work\" go"

peer0="$go run $ROOT/trab2/peer/main.go -p 8081 -n localhost:8083"
peer1="$go run $ROOT/trab2/peer/main.go -p 8082 -n localhost:8083"
peer2="$go run $ROOT/trab2/peer/main.go -p 8083 -n localhost:8084"
peer3="$go run $ROOT/trab2/peer/main.go -p 8084 -n localhost:8085"
peer5="$go run $ROOT/trab2/peer/main.go -p 8085"
peer4="$go run $ROOT/trab2/peer/main.go -p 8086 -n localhost:8084"

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
tmux send-keys -t 0 "$set_go $peer0" Enter
tmux send-keys -t 1 "$set_go $peer1" Enter
tmux send-keys -t 2 "$set_go $peer2" Enter
tmux send-keys -t 3 "$set_go $peer3" Enter
tmux send-keys -t 4 "$set_go $peer4" Enter
tmux send-keys -t 5 "$set_go $peer5" Enter
