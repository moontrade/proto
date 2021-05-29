# Moontrade Proto (MoonProto)

Moontrade heavily utilizes WebAssembly (AssemblyScript currently) for running custom code/algorithms to reduce pricing streams into custom data models, execute trades, re-balance portfolios, import/export transformations, etc.

Proto is optimized for flattened buffers and zero/minimal parsing. Given the heavy utilization of numerical computation, numerical types are stored in little-endian binary representation to natively support WASM without needing to transform/parse while supporting all native primitive types of WASM.

General purpose format. Even though WASM is the target, it can still compile schemas to other languages and runtimes (TypeScript/JavaScript, Go, Rust, C#, Python, etc) utilizing the lowest level features of the respective target as much as possible for peak performance and correctness.

# Schemas

Schemas are represented in ".moon" files. It's a bit of a hybrid between protobuf and flatbuffer schemas.

# Optimized for throughput

MoonProto is all about throughput over size. MoonProto messages can be compressed with a high-performance algorithm like LZ4 which Moontrade uses extensively.

# Optimized for streaming

MoonProto is built primarily with streaming in mind. Moontrade utilizes various streaming technologies like Redis, Kafka, PubSub, etc.

# Little-endian

MoonProto exclusively uses little-endian. Language implementations of MoonProto may provide support for big-endian format.

# Browser support

MoonProto supports TypeScript/JavaScript just dandy. Graphically charting numerical streams (financial charts) is extremely common in web browsers. Moontrade UI uses TypeScript HTML Canvas charting built on top of TradingView.

# Why not FlatBuffers?

FlatBuffers generally introduces an additional layer of indirection which can make property accesses slower with the added benefit of a more flexible evolutionary data structure. MoonProto still provides the ability to evolve your schemas, however MoonProto is geared for "native struct-like" access. In addition, it's important to point out that MoonProto was designed explicitly to meet Moontrade's product requirements / use cases. A good example of performance issues with FlatBuffers was evident inside Deno before they switched to their own direct (struct-like zero-copy) format. Similar performance gains can be expected for MoonProto.
