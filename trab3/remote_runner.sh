#!/usr/bin/env bash

work_dir="~/Dev/Go2Go/trab3"


LM2="go run peer/main.go L802.alunos.dcc.fc.up.pt"
LM3="go run peer/main.go L803.alunos.dcc.fc.up.pt"
LM4="go run peer/main.go L804.alunos.dcc.fc.up.pt"
LM5="go run peer/main.go L805.alunos.dcc.fc.up.pt"
LM6="go run peer/main.go L806.alunos.dcc.fc.up.pt"
LM7="go run peer/main.go L807.alunos.dcc.fc.up.pt"


# Send the commands to each of the 3 panes
tmux send-keys -t 2.0 "cd $work_dir; $LM2" Enter
tmux send-keys -t 2.1 "cd $work_dir; $LM3" Enter
tmux send-keys -t 2.2 "cd $work_dir; $LM4" Enter
tmux send-keys -t 2.3 "cd $work_dir; $LM5" Enter
tmux send-keys -t 2.4 "cd $work_dir; $LM6" Enter
tmux send-keys -t 2.5 "cd $work_dir; $LM7" Enter

tmux select-window -t 2
