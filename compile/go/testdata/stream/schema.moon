enum StreamKind : byte {
	Log = 0
	TimeSeries = 1
	Table = 2
}

enum SchemaKind : byte {
	Bytes			= 0		// Raw bytes
	MoonStruct		= 1		// MoonBuf Structure
	MoonMessage		= 1		// MoonBuf Message
	ProtoBuf		= 2		// Protocol buffers
	FlatBuf			= 3		// FlatBuffers
	Json			= 4		// Json
	MessagePack		= 5		// MessagePack
}

struct Stream {
	1	id			i64			// StreamID
	2	created		i64			// Unix timestamp of creation in nanoseconds
	3	accountID	i64			// AccountID that owns the stream
	4	duration	i64			// Duration of a single record. Only used if kind == Series
	5	name		string32	// Optional name
	6	record 		i32			// Record size
	7	kind		StreamKind	// Kind of stream
	8	schema		SchemaKind	// Schema serialization format
	9	realTime	bool		// Stream is appended in real-time
	10	blockSize	byte		// Size of default blocks (1, 2, 4, 8, 16, 32, 64)
}

struct AccountStats {
	id			i64
	storage		Stats
	appender	Stats
	streams		i64
}

struct StreamStats {
	storage		Stats			// Storage stats
	appender	Stats			// Appender stats
}

struct Stats {
	size		i64
	count		i64
	blocks		i64
}

// BlockID represents a globally unique ID of a single page of a single stream.
// String representation
struct BlockID {
	streamID	i64  	// StreamID
	id			i64		// Block ID / sequence
}

enum Compression : byte {
	None = 0
	LZ4 = 1
}

// BlockHeader
struct BlockHeader {
	streamID	i64				// Stream ID
	id			i64				// Block ID / Seq
	created		i64				// Unix Timestamp of creation in nanoseconds
	completed	i64				// Unix Timestamp of completion in nanoseconds
	min			i64				// Min record ID
	max			i64				// Max record ID
	start		i64				// Min timestamp
	end			i64				// Max timestamp
	savepoint	i64				// Current savepoint Block ID
	count		i32				// Number of records
	seq			i32				// Sequence number of first record
	size		i32				// Size of current data buffer
	sizeU		i32				// Size of data when uncompressed
	sizeX		i32				// Size of data when compressed
	compression	Compression		// Compression algorithm used
}

// Block64
struct Block64 {
	head		BlockHeader		// Header
	body		bytes65456		// Data
}

// Block32
struct Block32 {
	head		BlockHeader		// Header
	body		bytes32688		// Data
}

// Block16
struct Block16 {
	head		BlockHeader		// Header
	body		bytes16306		// Data
}

// Block8
struct Block8 {
	head		BlockHeader		// Header
	body		bytes8112		// Data
}

// Block4
struct Block4 {
	head		BlockHeader		// Header
	body		bytes4016		// Data
}

// Block2
struct Block2 {
	head		BlockHeader		// Header
	body		bytes1968		// Data
}

// Block1
struct Block1 {
	head		BlockHeader		// Header
	body		bytes944		// Data
}

struct RecordID {
	streamID	i64
	blockID		i64
	id			i64
}

enum MessageType : byte {
	Record 			= 1
	Block 			= 2
	EOS 			= 3
	EOB 			= 4
	Savepoint 		= 5
	Starting 		= 6
	Progress		= 7
	Started 		= 8
	Stopped 		= 9
}

enum StopReason : byte {
	// Stream is composed from another stream or external datasource and it stopped
	Source 		= 1
	// Stream has been paused
	Paused 		= 2
	// Stream is being migrated to a new writer
	Migrate 	= 3
	// Stream has stopped unexpectedly
	Error	 	= 4
}

struct RecordHeader {
	streamID	i64
	blockID		i64
	id			i64
	timestamp	i64
	start		i64
	end			i64
	savepoint	i64
	savepointR	i64
	seq			u16
	size		u16
	sizeU		u16
	sizeX		u16
	compression	Compression
	eob			bool
}

struct Savepoint {
	recordID		RecordID
	timestamp		i64
	writerID		i64			// ID of current writer that is appending the stream
}

// End of Stream
// The reader is caught up on the stream.
struct EOS {
	recordID	RecordID
	timestamp	i64
	writerID	i64			// ID of current writer that is appending the stream
	closed		bool
	waiting		bool
}

// End of Block
struct EOB {
	recordID	RecordID
	timestamp	i64
	savepoint	i64
}

struct Starting {
	recordID	RecordID	// Max record ID
	timestamp	i64			// Unix timestamp when message was created
	writerID	i64			// ID of current writer that is appending the stream
}

struct Progress {
	recordID	RecordID
	timestamp	i64
	writerID	i64			// ID of current writer that is appending the stream
	started		i64
	count		i64
	remaining	i64
}

struct Started {
	recordID	RecordID	// Max record ID
	timestamp	i64			// Unix timestamp when message was created
	writerID	i64			// ID of current writer that is appending the stream
	stops		i64			// Unix timestamp when stream will have a planned stop
}

struct Stopped {
    recordID		RecordID
	timestamp		i64				// Unix timestamp when message was created
	starts			i64				// Unix timestamp when stream is expected to start again
	reason			StopReason		// Reason stream was stopped
}
