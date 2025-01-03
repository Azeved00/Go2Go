#!/usr/bin/env bash


work_dir="~/Dev/Go2Go/trab2"

peer1="go run $work_dir/peer/main.go -n L803.alunos.dcc.fc.up.pt:8180"
peer2="go run $work_dir/peer/main.go -n L805.alunos.dcc.fc.up.pt:8180"
peer3="go run $work_dir/peer/main.go -n L803.alunos.dcc.fc.up.pt:8180"
peer4="go run $work_dir/peer/main.go -n L806.alunos.dcc.fc.up.pt:8180"
peer5="go run $work_dir/peer/main.go -n L805.alunos.dcc.fc.up.pt:8180"
peer6="go run $work_dir/peer/main.go "

# Send the commands to each of the 3 panes
tmux send-keys -t 2.0 "$peer1" Enter
tmux send-keys -t 2.1 "$peer2" Enter
tmux send-keys -t 2.2 "$peer3" Enter
tmux send-keys -t 2.3 "$peer4" Enter
tmux send-keys -t 2.4 "$peer5" Enter
tmux send-keys -t 2.5 "$peer6" Enter
