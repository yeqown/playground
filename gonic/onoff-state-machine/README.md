## Online or offline state machine

to judge online or offline state by the sequence of event. event includes: 
`enter` and `leave`. state includes `online` and `offline`.

For example (sequences):

* `E1 L1 E2 L2` as OFFLINE
* `E1 L2 L1 E1` as OFFLINE
* `E1 L1 L2 E2` as OFFLINE
* `E1 E2 E3 L1 L2` as ONLINE
* `E1 E2 L2 E3 L1` as ONLINE
* `E1 E2 L2 E3 L3 L1` as OFFLINE
* `E1 E2 L2 E3 L1 L3` as OFFLINE