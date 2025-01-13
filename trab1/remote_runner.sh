#!/usr/bin/env bash


work_dir="~/Dev/Go2Go/trab1"

LM2="go run $work_dir/server/main.go"

LM3="go run $work_dir/peer/main.go -s L1202.alunos.dcc.fc.up.pt:8080 L1204.alunos.dcc.fc.up.pt"
LM4="go run $work_dir/peer/main.go -s L1202.alunos.dcc.fc.up.pt:8080 L1205.alunos.dcc.fc.up.pt"
LM5="go run $work_dir/peer/main.go -s L1202.alunos.dcc.fc.up.pt:8080 L1206.alunos.dcc.fc.up.pt"
LM6="go run $work_dir/peer/main.go -s L1202.alunos.dcc.fc.up.pt:8080 L1207.alunos.dcc.fc.up.pt"
LM7="go run $work_dir/peer/main.go -s L1202.alunos.dcc.fc.up.pt:8080 L1203.alunos.dcc.fc.up.pt"

LM16="go run $work_dir/injector/main.go L1203.alunos.dcc.fc.up.pt:8180"

# Send the commands to each of the 3 panes
tmux send-keys -t 2.0 "$LM2" Enter
tmux send-keys -t 2.1 "$LM3" Enter
tmux send-keys -t 2.2 "$LM4" Enter
tmux send-keys -t 2.3 "$LM5" Enter
tmux send-keys -t 2.4 "$LM6" Enter
tmux send-keys -t 2.5 "$LM7" Enter
tmux select-window -t 2

sleep 1s
tmux send-keys -t 1.0 "$LM16" Enter

