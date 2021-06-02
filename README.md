# caps
CAPS (Continuous Automated Packet Sender) is a Go application that continuously sends packets to targets that traverse different networks.  The intent is to deploy CAPS on multiple Points of Presence, creating a mesh and understanding latency within the network.  This accounts for ingress/egress on anycast, as well, which is helpful for understanding first/last mile.

We store each target (per destination and network) in a buffer to process the average latency once the buffer is full.  This will bring awareness in the comparisson of the networks and provides understanding of which network should, or should not, be used.

Targets are obtained through an api; each target contains a network's ip:port to send to.  Each packet that is sent will be marked with a UUID to ensure the integrity of the packet: this will provide accurate latencies while traversing over multiple hops, as well as lowering potential risks of malformed packets.



Flow:

- Register the server as a target (it will not target itself)
- GET targets from the api
- On a timer, send a packet to each target
  - The packet is built and timestamped just before being sent
- When the Reflector on the server receives a packet, send it back to where it came from (changing the bool)
- When the Listener on the server receives a packet, timestamp and send to packet channel
- Create buffers for each destination and network
  - If it's already created, add to the latency array
- Once the buffer is full, send off for processing and reporting
  - An API is going to handle this and use Prometheus for metrics on latencies 
