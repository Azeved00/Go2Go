#  A Basic Chat application using Totally Ordered Multicast

The peers know all of the other peers and will, periodically, broadcast messages.
These Messages include a word (selected at random from a set) and a Lamport clock.

The Lamport clock includes a counter and the address of the machine that broadcasted the mssage.
These clocks are used to order the messages received.

## How to use?

Takes an address without port as a parameter (`0.0.0.0`),
which is the address of the machine this program is being run on.

Furthermore, there is an optional flag, 
`p`, port where the peer will listen, defaults to `8180`,

## Notes

The program is set up to turn it self off after 3 minutes.

