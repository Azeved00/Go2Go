# Counting the number of Nodes in a P2P Network

Using a data dissemination algorithm to count the number of nodes in a P2P network.

Each peer will keep a list of peers that they know of, 
the number of elements in the list will, eventually, approximate the total number of peers in the network.

In this example, the peers use a gossip algorithm to propagate their maps
which chooses, at random, a connection to gossip. Because of this, 
outer nodes of the network will take longer to know other outer nodes while 
more central nodes will have a more updated knowledge.

## How to use?

Takes an address without port as a parameter (`0.0.0.0`),
which is the address of the machine this program is being run on.

Furthermore, there are a few optional flags:

- `p` port where the peer will listen, defaults to `8180`,
- `n` address to which this peer will conect to, 
if no address is given then the peer will only wait for connections,


## Notes

The program is set up to turn it self off after 3 minutes.

