# Mutual Exclusion with Token Ring Algorithm


This algorithm connects a series of peers in a ring (1 -> 2 -> 3 ->1), 
after connection is established an injector is used to inject the token into one of the peers.

The token will be passed from peer to peer and only the person with the token can make requests to the server.

In this concrete example, while the peers do not have the token they will populate a request queue
which, when the peer gets the token, will be sent to the server.

## How to use?

There are 3 different modules included in this folder:

The [injector](./injector), the [peer](./peer), and the [server](server).

#### The injector

Takes an full address (`0.0.0.0:0000`) as a parameter. 
Connects to the peer in the address and injects the token.

#### The server

Will listen on port `8080` to the peer connections.
If no peers are connected for 30 seconds the server will close automatically.


#### The peer

Takes an address without port as a parameter (`0.0.0.0`),
which is the address of the machine that is next in the ring.

Furthermore, there are a few optional flags:

- `p` port where the peer will listen, defaults to `8180`,
- `n` port to which the peer will connect to, defaults to `8180`,
- `s` the address of the server, defaults to `localhost:8080`

The peers are set up to exit after 2 minutes of execution, furthermore, 
if a peer detects that the previous connection lo longer exists 
it will wait 30 seconds for them to reconnect, after which it will exit.


## Notes

The program is set up to turn it self off after 2 minutes. Furthermore, if a peer
