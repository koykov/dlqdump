# DLQ dump

Dumping DLQ implementation for [queue](https://github.com/koykov/queue) package. Allows dump leaky items to some storage:
disk, cloud, etc...

`DLQ` is an auxiliary queue, where drops leaky items from main queue due to workers limit reached (or error occurs).

This package also describes the opposite interface, that allows to restore items from the dump and return them back to
the main queue (see "Restoring" chapter).

## Dumping

Dumping-component implemented as queue (because `DLQ` must implement the corresponding
[interface]https://github.com/koykov/queue/blob/master/interface.go#L4 of the queue) and calls `Queue`.

The queue implements using [config](https://github.com/koykov/dlqdump/blob/master/config.go#L18).

The base param is a `Version`. This params will check on reading the dump and specify are data backward-compatibility or
not.

### Flush settings

Params `Capacity` and `FlushInterval` specifies how often the dumping queue must flush collected data to the dump.
`Capacity` must be set up in bytes; `FlushInterval` limits how long the queue must wait before flushing (beginning with
the moment of coming the first item in DLQ).

> **_NOTE:_**  `Capacity` must specify in bytes because dumping assumes the storing in some storage with limits the size
> (eg. some cloud with limit of file size). Thus, DLQ may collect limited amount of serialized data and will flush them
> by limit reach. This param is mandatory.
> 
> Similar work `FlushInterval`. It remembers the moment of coming the first item and by reaching `FlushInterval` time
> flushes the data with reason "interval reached". That param requires for cases when items comes to DLQ rarely and
> couldn't fill DLQ to `Capacity` limit. Because of `FlushInterval` the items will not stores in DLQ infinitely and will
> flush to storage even if size will small.

As a result, DLQ wait for coming the items and then checks what will occurs first: size of collected data will overflow
`Capacity` or `FlushInterval` will reach.

Note, on close the DLQ, the force flush will happen, independent of both params. Then DLQ will close.

### Serialization

Queue has two abstraction layers. The first is param `Encoder` - special component, that must implement interface
[`Encoder`](encoder.go). That component takes the arbitrary item and tries to serialize it to buffer `dst`.
Serialized data will send to the storage afterward.

`dlqdump` has several builtin encoders:
[builtin](encoder/builtin.go) и
[marshaller](encoder/marshaller.go).
The first may serialize string/bytes data or types implements `Byter` and `Stringer` interfaces.
The second may serialize [protobuf](https://en.wikipedia.org/wiki/Protocol_Buffers) objects.

### Dump writing

The second abstraction layer. There is a param `Writer` that must implement [`Writer`](writer.go) interface. This object,
using version and serialized data, writes a dump. `dlqdump` has builtin `Writer` implementation to 
[write dumps to the disk](fs).

You may write your own implementation to write dumps to the cloud, etc...

## Restoring

Dump writing isn't the full issues. Data from dumps should be used (restored and processed again). `dlqdump` contains
`Restorer` component, that are opposite to `Queue`.

The main idea: the source queue leaks and using DLQ sends the items to dumping queue. The queue flushed the data to
storage and then `Restorer` checks its periodically and tries to send items back to target queue (the origin queue in
most usable case, but you may specify any other queue).
As result, the loop is formed:
* queue leaks the items
* DLQ writes a dump
* Restorer reads the items from dump
* Restorer send restored items back to the queue

The storage uses as big buffer in that case, but not in RAM.

`Restorer` uses the same config struct, but ignores specific for queue params (queue similarly ignores `Restorer` params).

The base param is `Version`. Work similar to queue config. If version in config and dump will different, then dump will
be removed.

The target queue set up using param `Queue` and must implement [queue interface](https://github.com/koykov/queue/blob/master/interface.go#L4).

### Restoring settings

`Restorer` has three params:
* `CheckInterval` - the interval between checks of dumps in storage
* `PostponeInterval` - how long restoring must be postponed if target's queue rate overflows `AllowRate`
* `AllowRate` - the maximum rate (items/capacity) of target queue that allows to send items to it. Required to avoid
overflowing of target queue by `Restorer`

### Dump reading

`Restorer` similar to `Queue` has two abstraction layers, but in reverse meaning.

The first layer represents by param `Reader` that must implement [`Reader`](reader.go) interface. This object must read
from the dump version and serialized data till EOF error caught.

`dlqdump` has builtin implementation that [reads dump from the disk](fs). As usual, you may write your own implementation
for required storage.

### Deserialization

Serialized data taken from `Reader` will send to `Decoder` afterward - special param that must implement
[`Decoder`](decoder.go) interface. This object will deserialize the data or report about error occurs.

`dlqdump` has two builtin decoders:
[fallthrough](decoder/fallthrough.go) и
[unmarshaller](decoder/unmarshaller.go). The first one uses only for testing purposes. The second is opposite to
`marshaller` encoder and may deserialize objects like protobuf.

After success deserialization the item will send to the target queue.

