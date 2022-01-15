# NATS Logger

NATS uses a simple publish/subscribe style plain-text protocol to communicate between a NATS Server and its clients. Whilst this connection should remain opaque to the user, it can be quite handy to see the data being passed from time to time - this tool does just that (it also saves me loading Wireshark and filtering the NATS traffic). 

More information on NATS connection protocol can be found [here](https://docs.nats.io/reference/reference-protocols/nats-protocol).
    
<img width="1338" alt="screenshot" src="https://user-images.githubusercontent.com/1237341/149612435-6cdfac10-44ac-4bd9-8e4f-d8464893a235.png">

## Next step...

Next step is to add support for printing a pretty print version to a web UI over websockets.
