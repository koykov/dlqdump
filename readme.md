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

