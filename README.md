# caps
CAPS (Continuous Automated Packet Sender) is a Go application that continuously sends packets to compare multiple networks.  Each target will be stored in a buffer and latency results will be sent to an external application for metric processing.

Targets are obtained through an api and are continuously iterated over.  Each packet that is sent will be marked with a UUID to ensure the integrity of that packet that was sent: this will ensure accurate latencies while traversing over multiple hops, as well as lowering potential risks of malformed packets.

This binary requires:

An available ipv4 address to send, receive and reflect from
An available port to send, receive and reflect from

