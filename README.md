# WAP (WebAssembly Proto Format)

Moontrade heavily utilizes WebAssembly (TinyGo currently) for running custom code/algorithms to reduce pricing streams
into custom data models, execute trades, re-balance portfolios, import/export transformations, etc.

WAP is a flattened buffer format with zero parsing. Given the heavy utilization of numerical computation, numerical
types are stored in little-endian binary representation to natively support WASM without needing to transform/parse
while supporting all native primitive types of WASM. WAP supports other serialization formats which serialize and
deserialize into the flat format.

Given the low-level nature of Wasm, there doesn't exist a uniform way to describe data schemas efficiently without the
need to deserialze.

General purpose format. Even though WASM is a primary target, it has first class support to compile schemas to other
languages and runtimes (TypeScript/JavaScript, Go, Rust, C#, Python, etc) utilizing the lowest level features of the
respective target as much as possible for peak performance and correctness.

# Schemas

Schemas are represented in ".wap" files. It's a bit of a hybrid between protobuf and flatbuffer schemas.

# Optimized for throughput

MoonProto is all about throughput over size. MoonProto messages can be compressed with a high-performance algorithm like
LZ4 which Moontrade uses extensively.

# Optimized for streaming

MoonProto is built primarily with streaming in mind. Moontrade utilizes various streaming technologies like Redis,
Kafka, PubSub, etc.

# Little-endian

MoonProto exclusively uses little-endian. Language implementations of MoonProto may provide support for big-endian
format.

# JSON Support Built-in

WAP can transparently serialize and deserialize utilizing JSON format while taking advantage of the efficient flat
deserialized format.

### JSON Short and Long Names

WAP provides the ability to describe 2 names for the same field. Both must be unique within the record. The
auto-generated JSON reader will resolve either name.

# Protocol Buffers Support Built-in

WAP can transparently serialize and deserialize utilizing Protocol Buffers format while taking advantage of the
efficient flat deserialized format.

# Browser support

WAP supports TypeScript/JavaScript just dandy. Graphically charting numerical streams (financial charts) is extremely
common in web browsers. Moontrade UI uses TypeScript HTML Canvas charting built on top of TradingView.

# Variable Sized Buffers

For maximum performance structs are ideal. However, some structures require variable length strings, lists, etc. WAP can
support buffers up to 2GB in length.

# References

WAP supports references within a single graph. Adding a reference to a type will turn the type into a flexible length
type in order to resolve references in the graph.

# Anything can be a Root

# Reflection

# Streaming

WAP provides a custom streaming format.

### Logs

Logs timestamped monotonic records

### Series

Time-Series records are timestamped monotonic records and have a time range / intervals.

# Why not FlatBuffers?

FlatBuffers generally introduces an additional layer of indirection which can make property accesses slower with the
added benefit of a more flexible evolutionary data structure. WAP still provides the ability to evolve your schemas,
however WAP is geared for "native struct-like" access. In addition, it's important to point out that WAP was designed
explicitly to meet Moontrade's product requirements / use cases. A good example of performance issues with FlatBuffers
was evident inside Deno before they switched to their own direct (struct-like zero-copy) format. Similar performance
gains can be expected for WAP.

In addition, the mutable access pattern for FlatBuffers is awkward.
