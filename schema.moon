enum Kind : byte {
    Unknown         = 0
    Bool            = 1
    Byte            = 2
    Int8            = 3
    UInt8           = 4
    Int16           = 5
    UInt16          = 6
    Int32           = 7
    UInt32          = 8
    Int64           = 9
    UInt64          = 10
    Float32         = 11
    Float64         = 12
    String          = 13
    StringFixed     = 14
    Bytes           = 15
    BytesFixed      = 16
    RecordHeader    = 17
    BlockHeader     = 18
    Enum            = 30
    Record          = 40
    Struct          = 41
    List            = 50
}

enum BlockSize : u16 {
    B1kb  = 1024
    B2kb  = 2048
    B4kb  = 4096
    B8kb  = 8192
    B16kb = 16384
    B32kb = 32768
    B64kb = 65535
}

enum Format : byte {
    Raw         = 0
    Json        = 1
    Moon        = 2
    Protobuf    = 3
    Msgpack     = 4
}

enum Encoding : byte {
    None        = 0
    LZ4         = 1
    ZSTD        = 2
    Brotli      = 3
    Gzip        = 4
}

enum RecordLayout : byte {
    Aligned = 0
    Compact = 1
}

enum StreamKind : byte {
	Log     = 0
	Series  = 1
	Table   = 2
}

enum BlockLayout : byte {
    Row     = 1
    Column  = 2
}

// BlockHeader
struct BlockHeader {
    streamID    | a     i64
    id          | b     i64
    headID      | c     i64
    headMin     | d     i64
    headStart   | e     i64
    blocks      | f     i64
    records     | g     i64
	storage		| h     u64				// Cumulative storage usage including this block
	storageU	| i     u64				// Cumulative storage usage including this block when uncompressed
	created		| j     i64				// Unix Timestamp of creation in nanoseconds
	completed	| k     i64				// Unix Timestamp of completion in nanoseconds
	start		| l     i64				// Min timestamp
	end			| m     i64				// Max timestamp
	min			| n     i64				// Min record ID
	max			| o     i64				// Max record ID
	count		| p     u16				// Number of records
	size		| q     u16				// Size of current data buffer
	sizeU		| r     u16				// Size of data when uncompressed
	sizeX		| s     u16				// Size of data when compressed
	record		| t     u16				// Size of record. 0 = variable length
	blockSize	| u     BlockSize		// Size of uncompressed block (1024, 2048, 4096, 8192, 16384, 32768, 65535)
	encoding	| v     Compression		// Compression algorithm used
	kind		| w     StreamKind		// Kind of Stream (e.g. Log, Time-Series, or Table)
	format		| x     Format  		// Kind of serialization format
}

struct Stream {
	id			| a     i64			// StreamID
	created		| b     i64			// Unix timestamp of creation in nanoseconds
	accountID	| c     i64			// AccountID that owns the stream
	duration	| d     i64			// Duration of a single record. Only used if kind == Series
	record 		| e     u16			// Record size for fixed sized records
	name		| f     string32	// Optional name
	kind		| h     StreamKind	// Kind of stream
	format		| i     Format  	// Schema serialization format
	blockSize	| u     BlockSize	// Size of uncompressed block (1024, 2048, 4096, 8192, 16384, 32768, 65535)
}

// BlockID represents a globally unique ID of a single page of a single stream.
// String representation
struct BlockID {
	streamID	i64  	// StreamID
	id			i64		// Block ID / sequence
}

struct RecordID {
	streamID	i64
	blockID		i64
	id			i64
}

// EOS = End of Stream
// The reader is caught up on the stream and is NOT subscribed
// to new records.
struct EOS {
	recordID	RecordID
	timestamp	i64
}

// EOSWaiting = End of Stream Waiting for next record.
// The reader is caught up on the stream and is subscribed
// to new records.
struct EOSWaiting {
	recordID	RecordID
	timestamp	i64
}

enum MessageType : byte {
	Record 			= 1
	Records 		= 2
	Block 			= 3
	EOS 			= 4
	EOSWaiting 		= 5
	Savepoint 		= 6
	Starting 		= 7
	Started 		= 8
	Stopped 		= 9
	SyncStarted 	= 10
	SyncProgress 	= 11
	SyncStopped 	= 12
}

enum StopReason : byte {
	// Stream is composed from another stream or external datasource and it stopped
	Source 		= 1
	// Stream has been paused
	Paused 		= 2
	// Stream is being migrated to a new writer
	Migrate 	= 3
	// Stream has stopped unexpectedly
	Unexpected 	= 4
}

struct RecordHeader {
	id          RecordID
	prevID		i64
	timestamp	i64
	start		i64
	end			i64
	seq			u16
	size		u16
	sizeU		u16
	encoding	Encoding
	pad			bool
}

struct RecordsHeader {
	header 	RecordHeader
	count 	u16
	record 	u16
}

struct Savepoint {
	recordID		RecordID
	timestamp		i64
	duration		i64
}

struct SyncStarted {
	recordID	RecordID
	timestamp	i64
}

struct SyncProgress {
	recordID	RecordID
	timestamp	i64
	started		i64
	count		i64
	remaining	i64
}

enum SyncStoppedReason : byte {
	Success = 1
	Error = 2
}

struct SyncStopped {
	progress 	SyncProgress
	message		string64
	reason		SyncStoppedReason
}

struct Starting {
	recordID	RecordID	// Max record ID
	timestamp	i64			// Unix timestamp when message was created
	writerID	i64			// ID of current writer that is appending the stream
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
