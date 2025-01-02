# Simple Peer-To-Peer algorithms in Go

These are a collection of 3 p2p algorithms, implemented in Go.
The requirements for these implementations are described in the [assigment file](./assigment.pdf).

## The Algorithms

The first algorithm is a Mutual Exclusion with the Token Ring Algorithm, 
the second is to Count The number of Nodes in a P2P Network and the last 
a Basic Chat Application using Totally-Ordered Multicast, 
these are, respectivelly, in the folders [trab1](./trab1) [trab2](./trab2) and [trab3](./trab3)

## How To Run

In each of the folders there is a script, `runner.sh`, that sets up the network and starts all the peers.
Do note that the script uses `tmux` to run the peers in different panes.
