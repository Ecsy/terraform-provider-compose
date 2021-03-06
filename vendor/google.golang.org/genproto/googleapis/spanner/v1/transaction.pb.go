// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/spanner/v1/transaction.proto

package spanner

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf2 "github.com/golang/protobuf/ptypes/duration"
import google_protobuf3 "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// # Transactions
//
//
// Each session can have at most one active transaction at a time. After the
// active transaction is completed, the session can immediately be
// re-used for the next transaction. It is not necessary to create a
// new session for each transaction.
//
// # Transaction Modes
//
// Cloud Spanner supports two transaction modes:
//
//   1. Locking read-write. This type of transaction is the only way
//      to write data into Cloud Spanner. These transactions rely on
//      pessimistic locking and, if necessary, two-phase commit.
//      Locking read-write transactions may abort, requiring the
//      application to retry.
//
//   2. Snapshot read-only. This transaction type provides guaranteed
//      consistency across several reads, but does not allow
//      writes. Snapshot read-only transactions can be configured to
//      read at timestamps in the past. Snapshot read-only
//      transactions do not need to be committed.
//
// For transactions that only read, snapshot read-only transactions
// provide simpler semantics and are almost always faster. In
// particular, read-only transactions do not take locks, so they do
// not conflict with read-write transactions. As a consequence of not
// taking locks, they also do not abort, so retry loops are not needed.
//
// Transactions may only read/write data in a single database. They
// may, however, read/write data in different tables within that
// database.
//
// ## Locking Read-Write Transactions
//
// Locking transactions may be used to atomically read-modify-write
// data anywhere in a database. This type of transaction is externally
// consistent.
//
// Clients should attempt to minimize the amount of time a transaction
// is active. Faster transactions commit with higher probability
// and cause less contention. Cloud Spanner attempts to keep read locks
// active as long as the transaction continues to do reads, and the
// transaction has not been terminated by
// [Commit][google.spanner.v1.Spanner.Commit] or
// [Rollback][google.spanner.v1.Spanner.Rollback].  Long periods of
// inactivity at the client may cause Cloud Spanner to release a
// transaction's locks and abort it.
//
// Reads performed within a transaction acquire locks on the data
// being read. Writes can only be done at commit time, after all reads
// have been completed.
// Conceptually, a read-write transaction consists of zero or more
// reads or SQL queries followed by
// [Commit][google.spanner.v1.Spanner.Commit]. At any time before
// [Commit][google.spanner.v1.Spanner.Commit], the client can send a
// [Rollback][google.spanner.v1.Spanner.Rollback] request to abort the
// transaction.
//
// ### Semantics
//
// Cloud Spanner can commit the transaction if all read locks it acquired
// are still valid at commit time, and it is able to acquire write
// locks for all writes. Cloud Spanner can abort the transaction for any
// reason. If a commit attempt returns `ABORTED`, Cloud Spanner guarantees
// that the transaction has not modified any user data in Cloud Spanner.
//
// Unless the transaction commits, Cloud Spanner makes no guarantees about
// how long the transaction's locks were held for. It is an error to
// use Cloud Spanner locks for any sort of mutual exclusion other than
// between Cloud Spanner transactions themselves.
//
// ### Retrying Aborted Transactions
//
// When a transaction aborts, the application can choose to retry the
// whole transaction again. To maximize the chances of successfully
// committing the retry, the client should execute the retry in the
// same session as the original attempt. The original session's lock
// priority increases with each consecutive abort, meaning that each
// attempt has a slightly better chance of success than the previous.
//
// Under some circumstances (e.g., many transactions attempting to
// modify the same row(s)), a transaction can abort many times in a
// short period before successfully committing. Thus, it is not a good
// idea to cap the number of retries a transaction can attempt;
// instead, it is better to limit the total amount of wall time spent
// retrying.
//
// ### Idle Transactions
//
// A transaction is considered idle if it has no outstanding reads or
// SQL queries and has not started a read or SQL query within the last 10
// seconds. Idle transactions can be aborted by Cloud Spanner so that they
// don't hold on to locks indefinitely. In that case, the commit will
// fail with error `ABORTED`.
//
// If this behavior is undesirable, periodically executing a simple
// SQL query in the transaction (e.g., `SELECT 1`) prevents the
// transaction from becoming idle.
//
// ## Snapshot Read-Only Transactions
//
// Snapshot read-only transactions provides a simpler method than
// locking read-write transactions for doing several consistent
// reads. However, this type of transaction does not support writes.
//
// Snapshot transactions do not take locks. Instead, they work by
// choosing a Cloud Spanner timestamp, then executing all reads at that
// timestamp. Since they do not acquire locks, they do not block
// concurrent read-write transactions.
//
// Unlike locking read-write transactions, snapshot read-only
// transactions never abort. They can fail if the chosen read
// timestamp is garbage collected; however, the default garbage
// collection policy is generous enough that most applications do not
// need to worry about this in practice.
//
// Snapshot read-only transactions do not need to call
// [Commit][google.spanner.v1.Spanner.Commit] or
// [Rollback][google.spanner.v1.Spanner.Rollback] (and in fact are not
// permitted to do so).
//
// To execute a snapshot transaction, the client specifies a timestamp
// bound, which tells Cloud Spanner how to choose a read timestamp.
//
// The types of timestamp bound are:
//
//   - Strong (the default).
//   - Bounded staleness.
//   - Exact staleness.
//
// If the Cloud Spanner database to be read is geographically distributed,
// stale read-only transactions can execute more quickly than strong
// or read-write transaction, because they are able to execute far
// from the leader replica.
//
// Each type of timestamp bound is discussed in detail below.
//
// ### Strong
//
// Strong reads are guaranteed to see the effects of all transactions
// that have committed before the start of the read. Furthermore, all
// rows yielded by a single read are consistent with each other -- if
// any part of the read observes a transaction, all parts of the read
// see the transaction.
//
// Strong reads are not repeatable: two consecutive strong read-only
// transactions might return inconsistent results if there are
// concurrent writes. If consistency across reads is required, the
// reads should be executed within a transaction or at an exact read
// timestamp.
//
// See [TransactionOptions.ReadOnly.strong][google.spanner.v1.TransactionOptions.ReadOnly.strong].
//
// ### Exact Staleness
//
// These timestamp bounds execute reads at a user-specified
// timestamp. Reads at a timestamp are guaranteed to see a consistent
// prefix of the global transaction history: they observe
// modifications done by all transactions with a commit timestamp <=
// the read timestamp, and observe none of the modifications done by
// transactions with a larger commit timestamp. They will block until
// all conflicting transactions that may be assigned commit timestamps
// <= the read timestamp have finished.
//
// The timestamp can either be expressed as an absolute Cloud Spanner commit
// timestamp or a staleness relative to the current time.
//
// These modes do not require a "negotiation phase" to pick a
// timestamp. As a result, they execute slightly faster than the
// equivalent boundedly stale concurrency modes. On the other hand,
// boundedly stale reads usually return fresher results.
//
// See [TransactionOptions.ReadOnly.read_timestamp][google.spanner.v1.TransactionOptions.ReadOnly.read_timestamp] and
// [TransactionOptions.ReadOnly.exact_staleness][google.spanner.v1.TransactionOptions.ReadOnly.exact_staleness].
//
// ### Bounded Staleness
//
// Bounded staleness modes allow Cloud Spanner to pick the read timestamp,
// subject to a user-provided staleness bound. Cloud Spanner chooses the
// newest timestamp within the staleness bound that allows execution
// of the reads at the closest available replica without blocking.
//
// All rows yielded are consistent with each other -- if any part of
// the read observes a transaction, all parts of the read see the
// transaction. Boundedly stale reads are not repeatable: two stale
// reads, even if they use the same staleness bound, can execute at
// different timestamps and thus return inconsistent results.
//
// Boundedly stale reads execute in two phases: the first phase
// negotiates a timestamp among all replicas needed to serve the
// read. In the second phase, reads are executed at the negotiated
// timestamp.
//
// As a result of the two phase execution, bounded staleness reads are
// usually a little slower than comparable exact staleness
// reads. However, they are typically able to return fresher
// results, and are more likely to execute at the closest replica.
//
// Because the timestamp negotiation requires up-front knowledge of
// which rows will be read, it can only be used with single-use
// read-only transactions.
//
// See [TransactionOptions.ReadOnly.max_staleness][google.spanner.v1.TransactionOptions.ReadOnly.max_staleness] and
// [TransactionOptions.ReadOnly.min_read_timestamp][google.spanner.v1.TransactionOptions.ReadOnly.min_read_timestamp].
//
// ### Old Read Timestamps and Garbage Collection
//
// Cloud Spanner continuously garbage collects deleted and overwritten data
// in the background to reclaim storage space. This process is known
// as "version GC". By default, version GC reclaims versions after they
// are one hour old. Because of this, Cloud Spanner cannot perform reads
// at read timestamps more than one hour in the past. This
// restriction also applies to in-progress reads and/or SQL queries whose
// timestamp become too old while executing. Reads and SQL queries with
// too-old read timestamps fail with the error `FAILED_PRECONDITION`.
type TransactionOptions struct {
	// Required. The type of transaction.
	//
	// Types that are valid to be assigned to Mode:
	//	*TransactionOptions_ReadWrite_
	//	*TransactionOptions_ReadOnly_
	Mode isTransactionOptions_Mode `protobuf_oneof:"mode"`
}

func (m *TransactionOptions) Reset()                    { *m = TransactionOptions{} }
func (m *TransactionOptions) String() string            { return proto.CompactTextString(m) }
func (*TransactionOptions) ProtoMessage()               {}
func (*TransactionOptions) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{0} }

type isTransactionOptions_Mode interface {
	isTransactionOptions_Mode()
}

type TransactionOptions_ReadWrite_ struct {
	ReadWrite *TransactionOptions_ReadWrite `protobuf:"bytes,1,opt,name=read_write,json=readWrite,oneof"`
}
type TransactionOptions_ReadOnly_ struct {
	ReadOnly *TransactionOptions_ReadOnly `protobuf:"bytes,2,opt,name=read_only,json=readOnly,oneof"`
}

func (*TransactionOptions_ReadWrite_) isTransactionOptions_Mode() {}
func (*TransactionOptions_ReadOnly_) isTransactionOptions_Mode()  {}

func (m *TransactionOptions) GetMode() isTransactionOptions_Mode {
	if m != nil {
		return m.Mode
	}
	return nil
}

func (m *TransactionOptions) GetReadWrite() *TransactionOptions_ReadWrite {
	if x, ok := m.GetMode().(*TransactionOptions_ReadWrite_); ok {
		return x.ReadWrite
	}
	return nil
}

func (m *TransactionOptions) GetReadOnly() *TransactionOptions_ReadOnly {
	if x, ok := m.GetMode().(*TransactionOptions_ReadOnly_); ok {
		return x.ReadOnly
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*TransactionOptions) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TransactionOptions_OneofMarshaler, _TransactionOptions_OneofUnmarshaler, _TransactionOptions_OneofSizer, []interface{}{
		(*TransactionOptions_ReadWrite_)(nil),
		(*TransactionOptions_ReadOnly_)(nil),
	}
}

func _TransactionOptions_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*TransactionOptions)
	// mode
	switch x := m.Mode.(type) {
	case *TransactionOptions_ReadWrite_:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ReadWrite); err != nil {
			return err
		}
	case *TransactionOptions_ReadOnly_:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ReadOnly); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("TransactionOptions.Mode has unexpected type %T", x)
	}
	return nil
}

func _TransactionOptions_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*TransactionOptions)
	switch tag {
	case 1: // mode.read_write
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TransactionOptions_ReadWrite)
		err := b.DecodeMessage(msg)
		m.Mode = &TransactionOptions_ReadWrite_{msg}
		return true, err
	case 2: // mode.read_only
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TransactionOptions_ReadOnly)
		err := b.DecodeMessage(msg)
		m.Mode = &TransactionOptions_ReadOnly_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _TransactionOptions_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*TransactionOptions)
	// mode
	switch x := m.Mode.(type) {
	case *TransactionOptions_ReadWrite_:
		s := proto.Size(x.ReadWrite)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *TransactionOptions_ReadOnly_:
		s := proto.Size(x.ReadOnly)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Message type to initiate a read-write transaction. Currently this
// transaction type has no options.
type TransactionOptions_ReadWrite struct {
}

func (m *TransactionOptions_ReadWrite) Reset()                    { *m = TransactionOptions_ReadWrite{} }
func (m *TransactionOptions_ReadWrite) String() string            { return proto.CompactTextString(m) }
func (*TransactionOptions_ReadWrite) ProtoMessage()               {}
func (*TransactionOptions_ReadWrite) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{0, 0} }

// Message type to initiate a read-only transaction.
type TransactionOptions_ReadOnly struct {
	// How to choose the timestamp for the read-only transaction.
	//
	// Types that are valid to be assigned to TimestampBound:
	//	*TransactionOptions_ReadOnly_Strong
	//	*TransactionOptions_ReadOnly_MinReadTimestamp
	//	*TransactionOptions_ReadOnly_MaxStaleness
	//	*TransactionOptions_ReadOnly_ReadTimestamp
	//	*TransactionOptions_ReadOnly_ExactStaleness
	TimestampBound isTransactionOptions_ReadOnly_TimestampBound `protobuf_oneof:"timestamp_bound"`
	// If true, the Cloud Spanner-selected read timestamp is included in
	// the [Transaction][google.spanner.v1.Transaction] message that describes the transaction.
	ReturnReadTimestamp bool `protobuf:"varint,6,opt,name=return_read_timestamp,json=returnReadTimestamp" json:"return_read_timestamp,omitempty"`
}

func (m *TransactionOptions_ReadOnly) Reset()                    { *m = TransactionOptions_ReadOnly{} }
func (m *TransactionOptions_ReadOnly) String() string            { return proto.CompactTextString(m) }
func (*TransactionOptions_ReadOnly) ProtoMessage()               {}
func (*TransactionOptions_ReadOnly) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{0, 1} }

type isTransactionOptions_ReadOnly_TimestampBound interface {
	isTransactionOptions_ReadOnly_TimestampBound()
}

type TransactionOptions_ReadOnly_Strong struct {
	Strong bool `protobuf:"varint,1,opt,name=strong,oneof"`
}
type TransactionOptions_ReadOnly_MinReadTimestamp struct {
	MinReadTimestamp *google_protobuf3.Timestamp `protobuf:"bytes,2,opt,name=min_read_timestamp,json=minReadTimestamp,oneof"`
}
type TransactionOptions_ReadOnly_MaxStaleness struct {
	MaxStaleness *google_protobuf2.Duration `protobuf:"bytes,3,opt,name=max_staleness,json=maxStaleness,oneof"`
}
type TransactionOptions_ReadOnly_ReadTimestamp struct {
	ReadTimestamp *google_protobuf3.Timestamp `protobuf:"bytes,4,opt,name=read_timestamp,json=readTimestamp,oneof"`
}
type TransactionOptions_ReadOnly_ExactStaleness struct {
	ExactStaleness *google_protobuf2.Duration `protobuf:"bytes,5,opt,name=exact_staleness,json=exactStaleness,oneof"`
}

func (*TransactionOptions_ReadOnly_Strong) isTransactionOptions_ReadOnly_TimestampBound()           {}
func (*TransactionOptions_ReadOnly_MinReadTimestamp) isTransactionOptions_ReadOnly_TimestampBound() {}
func (*TransactionOptions_ReadOnly_MaxStaleness) isTransactionOptions_ReadOnly_TimestampBound()     {}
func (*TransactionOptions_ReadOnly_ReadTimestamp) isTransactionOptions_ReadOnly_TimestampBound()    {}
func (*TransactionOptions_ReadOnly_ExactStaleness) isTransactionOptions_ReadOnly_TimestampBound()   {}

func (m *TransactionOptions_ReadOnly) GetTimestampBound() isTransactionOptions_ReadOnly_TimestampBound {
	if m != nil {
		return m.TimestampBound
	}
	return nil
}

func (m *TransactionOptions_ReadOnly) GetStrong() bool {
	if x, ok := m.GetTimestampBound().(*TransactionOptions_ReadOnly_Strong); ok {
		return x.Strong
	}
	return false
}

func (m *TransactionOptions_ReadOnly) GetMinReadTimestamp() *google_protobuf3.Timestamp {
	if x, ok := m.GetTimestampBound().(*TransactionOptions_ReadOnly_MinReadTimestamp); ok {
		return x.MinReadTimestamp
	}
	return nil
}

func (m *TransactionOptions_ReadOnly) GetMaxStaleness() *google_protobuf2.Duration {
	if x, ok := m.GetTimestampBound().(*TransactionOptions_ReadOnly_MaxStaleness); ok {
		return x.MaxStaleness
	}
	return nil
}

func (m *TransactionOptions_ReadOnly) GetReadTimestamp() *google_protobuf3.Timestamp {
	if x, ok := m.GetTimestampBound().(*TransactionOptions_ReadOnly_ReadTimestamp); ok {
		return x.ReadTimestamp
	}
	return nil
}

func (m *TransactionOptions_ReadOnly) GetExactStaleness() *google_protobuf2.Duration {
	if x, ok := m.GetTimestampBound().(*TransactionOptions_ReadOnly_ExactStaleness); ok {
		return x.ExactStaleness
	}
	return nil
}

func (m *TransactionOptions_ReadOnly) GetReturnReadTimestamp() bool {
	if m != nil {
		return m.ReturnReadTimestamp
	}
	return false
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*TransactionOptions_ReadOnly) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TransactionOptions_ReadOnly_OneofMarshaler, _TransactionOptions_ReadOnly_OneofUnmarshaler, _TransactionOptions_ReadOnly_OneofSizer, []interface{}{
		(*TransactionOptions_ReadOnly_Strong)(nil),
		(*TransactionOptions_ReadOnly_MinReadTimestamp)(nil),
		(*TransactionOptions_ReadOnly_MaxStaleness)(nil),
		(*TransactionOptions_ReadOnly_ReadTimestamp)(nil),
		(*TransactionOptions_ReadOnly_ExactStaleness)(nil),
	}
}

func _TransactionOptions_ReadOnly_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*TransactionOptions_ReadOnly)
	// timestamp_bound
	switch x := m.TimestampBound.(type) {
	case *TransactionOptions_ReadOnly_Strong:
		t := uint64(0)
		if x.Strong {
			t = 1
		}
		b.EncodeVarint(1<<3 | proto.WireVarint)
		b.EncodeVarint(t)
	case *TransactionOptions_ReadOnly_MinReadTimestamp:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MinReadTimestamp); err != nil {
			return err
		}
	case *TransactionOptions_ReadOnly_MaxStaleness:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MaxStaleness); err != nil {
			return err
		}
	case *TransactionOptions_ReadOnly_ReadTimestamp:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ReadTimestamp); err != nil {
			return err
		}
	case *TransactionOptions_ReadOnly_ExactStaleness:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExactStaleness); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("TransactionOptions_ReadOnly.TimestampBound has unexpected type %T", x)
	}
	return nil
}

func _TransactionOptions_ReadOnly_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*TransactionOptions_ReadOnly)
	switch tag {
	case 1: // timestamp_bound.strong
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.TimestampBound = &TransactionOptions_ReadOnly_Strong{x != 0}
		return true, err
	case 2: // timestamp_bound.min_read_timestamp
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf3.Timestamp)
		err := b.DecodeMessage(msg)
		m.TimestampBound = &TransactionOptions_ReadOnly_MinReadTimestamp{msg}
		return true, err
	case 3: // timestamp_bound.max_staleness
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf2.Duration)
		err := b.DecodeMessage(msg)
		m.TimestampBound = &TransactionOptions_ReadOnly_MaxStaleness{msg}
		return true, err
	case 4: // timestamp_bound.read_timestamp
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf3.Timestamp)
		err := b.DecodeMessage(msg)
		m.TimestampBound = &TransactionOptions_ReadOnly_ReadTimestamp{msg}
		return true, err
	case 5: // timestamp_bound.exact_staleness
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(google_protobuf2.Duration)
		err := b.DecodeMessage(msg)
		m.TimestampBound = &TransactionOptions_ReadOnly_ExactStaleness{msg}
		return true, err
	default:
		return false, nil
	}
}

func _TransactionOptions_ReadOnly_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*TransactionOptions_ReadOnly)
	// timestamp_bound
	switch x := m.TimestampBound.(type) {
	case *TransactionOptions_ReadOnly_Strong:
		n += proto.SizeVarint(1<<3 | proto.WireVarint)
		n += 1
	case *TransactionOptions_ReadOnly_MinReadTimestamp:
		s := proto.Size(x.MinReadTimestamp)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *TransactionOptions_ReadOnly_MaxStaleness:
		s := proto.Size(x.MaxStaleness)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *TransactionOptions_ReadOnly_ReadTimestamp:
		s := proto.Size(x.ReadTimestamp)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *TransactionOptions_ReadOnly_ExactStaleness:
		s := proto.Size(x.ExactStaleness)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// A transaction.
type Transaction struct {
	// `id` may be used to identify the transaction in subsequent
	// [Read][google.spanner.v1.Spanner.Read],
	// [ExecuteSql][google.spanner.v1.Spanner.ExecuteSql],
	// [Commit][google.spanner.v1.Spanner.Commit], or
	// [Rollback][google.spanner.v1.Spanner.Rollback] calls.
	//
	// Single-use read-only transactions do not have IDs, because
	// single-use transactions do not support multiple requests.
	Id []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// For snapshot read-only transactions, the read timestamp chosen
	// for the transaction. Not returned by default: see
	// [TransactionOptions.ReadOnly.return_read_timestamp][google.spanner.v1.TransactionOptions.ReadOnly.return_read_timestamp].
	//
	// A timestamp in RFC3339 UTC \"Zulu\" format, accurate to nanoseconds.
	// Example: `"2014-10-02T15:01:23.045123456Z"`.
	ReadTimestamp *google_protobuf3.Timestamp `protobuf:"bytes,2,opt,name=read_timestamp,json=readTimestamp" json:"read_timestamp,omitempty"`
}

func (m *Transaction) Reset()                    { *m = Transaction{} }
func (m *Transaction) String() string            { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()               {}
func (*Transaction) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{1} }

func (m *Transaction) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Transaction) GetReadTimestamp() *google_protobuf3.Timestamp {
	if m != nil {
		return m.ReadTimestamp
	}
	return nil
}

// This message is used to select the transaction in which a
// [Read][google.spanner.v1.Spanner.Read] or
// [ExecuteSql][google.spanner.v1.Spanner.ExecuteSql] call runs.
//
// See [TransactionOptions][google.spanner.v1.TransactionOptions] for more information about transactions.
type TransactionSelector struct {
	// If no fields are set, the default is a single use transaction
	// with strong concurrency.
	//
	// Types that are valid to be assigned to Selector:
	//	*TransactionSelector_SingleUse
	//	*TransactionSelector_Id
	//	*TransactionSelector_Begin
	Selector isTransactionSelector_Selector `protobuf_oneof:"selector"`
}

func (m *TransactionSelector) Reset()                    { *m = TransactionSelector{} }
func (m *TransactionSelector) String() string            { return proto.CompactTextString(m) }
func (*TransactionSelector) ProtoMessage()               {}
func (*TransactionSelector) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{2} }

type isTransactionSelector_Selector interface {
	isTransactionSelector_Selector()
}

type TransactionSelector_SingleUse struct {
	SingleUse *TransactionOptions `protobuf:"bytes,1,opt,name=single_use,json=singleUse,oneof"`
}
type TransactionSelector_Id struct {
	Id []byte `protobuf:"bytes,2,opt,name=id,proto3,oneof"`
}
type TransactionSelector_Begin struct {
	Begin *TransactionOptions `protobuf:"bytes,3,opt,name=begin,oneof"`
}

func (*TransactionSelector_SingleUse) isTransactionSelector_Selector() {}
func (*TransactionSelector_Id) isTransactionSelector_Selector()        {}
func (*TransactionSelector_Begin) isTransactionSelector_Selector()     {}

func (m *TransactionSelector) GetSelector() isTransactionSelector_Selector {
	if m != nil {
		return m.Selector
	}
	return nil
}

func (m *TransactionSelector) GetSingleUse() *TransactionOptions {
	if x, ok := m.GetSelector().(*TransactionSelector_SingleUse); ok {
		return x.SingleUse
	}
	return nil
}

func (m *TransactionSelector) GetId() []byte {
	if x, ok := m.GetSelector().(*TransactionSelector_Id); ok {
		return x.Id
	}
	return nil
}

func (m *TransactionSelector) GetBegin() *TransactionOptions {
	if x, ok := m.GetSelector().(*TransactionSelector_Begin); ok {
		return x.Begin
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*TransactionSelector) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TransactionSelector_OneofMarshaler, _TransactionSelector_OneofUnmarshaler, _TransactionSelector_OneofSizer, []interface{}{
		(*TransactionSelector_SingleUse)(nil),
		(*TransactionSelector_Id)(nil),
		(*TransactionSelector_Begin)(nil),
	}
}

func _TransactionSelector_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*TransactionSelector)
	// selector
	switch x := m.Selector.(type) {
	case *TransactionSelector_SingleUse:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SingleUse); err != nil {
			return err
		}
	case *TransactionSelector_Id:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeRawBytes(x.Id)
	case *TransactionSelector_Begin:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Begin); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("TransactionSelector.Selector has unexpected type %T", x)
	}
	return nil
}

func _TransactionSelector_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*TransactionSelector)
	switch tag {
	case 1: // selector.single_use
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TransactionOptions)
		err := b.DecodeMessage(msg)
		m.Selector = &TransactionSelector_SingleUse{msg}
		return true, err
	case 2: // selector.id
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeRawBytes(true)
		m.Selector = &TransactionSelector_Id{x}
		return true, err
	case 3: // selector.begin
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TransactionOptions)
		err := b.DecodeMessage(msg)
		m.Selector = &TransactionSelector_Begin{msg}
		return true, err
	default:
		return false, nil
	}
}

func _TransactionSelector_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*TransactionSelector)
	// selector
	switch x := m.Selector.(type) {
	case *TransactionSelector_SingleUse:
		s := proto.Size(x.SingleUse)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *TransactionSelector_Id:
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Id)))
		n += len(x.Id)
	case *TransactionSelector_Begin:
		s := proto.Size(x.Begin)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*TransactionOptions)(nil), "google.spanner.v1.TransactionOptions")
	proto.RegisterType((*TransactionOptions_ReadWrite)(nil), "google.spanner.v1.TransactionOptions.ReadWrite")
	proto.RegisterType((*TransactionOptions_ReadOnly)(nil), "google.spanner.v1.TransactionOptions.ReadOnly")
	proto.RegisterType((*Transaction)(nil), "google.spanner.v1.Transaction")
	proto.RegisterType((*TransactionSelector)(nil), "google.spanner.v1.TransactionSelector")
}

func init() { proto.RegisterFile("google/spanner/v1/transaction.proto", fileDescriptor5) }

var fileDescriptor5 = []byte{
	// 522 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x86, 0xd3, 0x6e, 0xab, 0xba, 0xd3, 0xae, 0xeb, 0x3c, 0x4d, 0x94, 0x08, 0x01, 0x2a, 0x42,
	0xe2, 0xca, 0x51, 0xc7, 0x0d, 0x12, 0x42, 0x82, 0x6e, 0x82, 0x08, 0x09, 0xad, 0x4a, 0x07, 0x48,
	0xdc, 0x04, 0xb7, 0x31, 0x91, 0xa5, 0xc4, 0x8e, 0x6c, 0x67, 0x74, 0x57, 0xdc, 0xf0, 0x34, 0x3c,
	0x02, 0x6f, 0xc1, 0x1b, 0xa1, 0x38, 0x4e, 0x9b, 0x35, 0x17, 0xeb, 0x5d, 0xdc, 0xf3, 0xff, 0xff,
	0xf9, 0x7c, 0x8e, 0x55, 0x78, 0x16, 0x0b, 0x11, 0x27, 0xd4, 0x53, 0x19, 0xe1, 0x9c, 0x4a, 0xef,
	0x66, 0xe2, 0x69, 0x49, 0xb8, 0x22, 0x4b, 0xcd, 0x04, 0xc7, 0x99, 0x14, 0x5a, 0xa0, 0x93, 0x52,
	0x84, 0xad, 0x08, 0xdf, 0x4c, 0xdc, 0x47, 0xd6, 0x47, 0x32, 0xe6, 0x11, 0xce, 0x85, 0x26, 0x85,
	0x5e, 0x95, 0x06, 0xf7, 0xb1, 0xad, 0x9a, 0xd3, 0x22, 0xff, 0xe1, 0x45, 0xb9, 0x24, 0x9b, 0x40,
	0xf7, 0xc9, 0x76, 0x5d, 0xb3, 0x94, 0x2a, 0x4d, 0xd2, 0xac, 0x14, 0x8c, 0xff, 0xed, 0x03, 0xba,
	0xde, 0x70, 0x5c, 0x65, 0x26, 0x1d, 0xcd, 0x00, 0x24, 0x25, 0x51, 0xf8, 0x53, 0x32, 0x4d, 0x47,
	0xad, 0xa7, 0xad, 0x17, 0xbd, 0x73, 0x0f, 0x37, 0xe8, 0x70, 0xd3, 0x8a, 0x03, 0x4a, 0xa2, 0xaf,
	0x85, 0xcd, 0x77, 0x82, 0x43, 0x59, 0x1d, 0xd0, 0x27, 0x30, 0x87, 0x50, 0xf0, 0xe4, 0x76, 0xd4,
	0x36, 0x81, 0x78, 0xf7, 0xc0, 0x2b, 0x9e, 0xdc, 0xfa, 0x4e, 0xd0, 0x95, 0xf6, 0xdb, 0xed, 0xc1,
	0xe1, 0xba, 0x91, 0xfb, 0x7b, 0x0f, 0xba, 0x95, 0x0a, 0x8d, 0xa0, 0xa3, 0xb4, 0x14, 0x3c, 0x36,
	0xd8, 0x5d, 0xdf, 0x09, 0xec, 0x19, 0x7d, 0x04, 0x94, 0x32, 0x1e, 0x1a, 0x8c, 0xf5, 0x1c, 0x2c,
	0x8b, 0x5b, 0xb1, 0x54, 0x93, 0xc2, 0xd7, 0x95, 0xc2, 0x77, 0x82, 0x61, 0xca, 0x78, 0xd1, 0x60,
	0xfd, 0x1b, 0x7a, 0x0b, 0x47, 0x29, 0x59, 0x85, 0x4a, 0x93, 0x84, 0x72, 0xaa, 0xd4, 0x68, 0xcf,
	0xc4, 0x3c, 0x6c, 0xc4, 0x5c, 0xda, 0x85, 0xf8, 0x4e, 0xd0, 0x4f, 0xc9, 0x6a, 0x5e, 0x19, 0xd0,
	0x05, 0x0c, 0xb6, 0x48, 0xf6, 0x77, 0x20, 0x39, 0x92, 0x77, 0x30, 0x2e, 0xe1, 0x98, 0xae, 0xc8,
	0x52, 0xd7, 0x40, 0x0e, 0xee, 0x07, 0x19, 0x18, 0xcf, 0x06, 0xe5, 0x1c, 0xce, 0x24, 0xd5, 0xb9,
	0x6c, 0xcc, 0xa6, 0x53, 0x4c, 0x30, 0x38, 0x2d, 0x8b, 0x77, 0x06, 0x30, 0x3d, 0x81, 0xe3, 0xb5,
	0x2e, 0x5c, 0x88, 0x9c, 0x47, 0xd3, 0x0e, 0xec, 0xa7, 0x22, 0xa2, 0xe3, 0xef, 0xd0, 0xab, 0xad,
	0x11, 0x0d, 0xa0, 0xcd, 0x22, 0xb3, 0x8c, 0x7e, 0xd0, 0x66, 0x11, 0x7a, 0xd7, 0xb8, 0xf8, 0xbd,
	0x2b, 0xd8, 0xba, 0xf6, 0xf8, 0x6f, 0x0b, 0x4e, 0x6b, 0x2d, 0xe6, 0x34, 0xa1, 0x4b, 0x2d, 0x24,
	0x7a, 0x0f, 0xa0, 0x18, 0x8f, 0x13, 0x1a, 0xe6, 0xaa, 0x7a, 0xb6, 0xcf, 0x77, 0x7a, 0x65, 0xc5,
	0x63, 0x2d, 0xad, 0x9f, 0x15, 0x45, 0x43, 0x83, 0x5c, 0x60, 0xf5, 0x7d, 0xc7, 0x40, 0xbf, 0x81,
	0x83, 0x05, 0x8d, 0x19, 0xb7, 0x7b, 0xde, 0x39, 0xb4, 0x74, 0x4d, 0x01, 0xba, 0xca, 0x42, 0x4e,
	0x7f, 0xc1, 0xd9, 0x52, 0xa4, 0xcd, 0x80, 0xe9, 0xb0, 0x96, 0x30, 0x2b, 0x66, 0x30, 0x6b, 0x7d,
	0x7b, 0x65, 0x65, 0xb1, 0x48, 0x08, 0x8f, 0xb1, 0x90, 0xb1, 0x17, 0x53, 0x6e, 0x26, 0xe4, 0x95,
	0x25, 0x92, 0x31, 0x55, 0xfb, 0x57, 0x79, 0x6d, 0x3f, 0xff, 0xb4, 0x1f, 0x7c, 0x28, 0xad, 0x17,
	0x89, 0xc8, 0x23, 0x3c, 0xb7, 0x7d, 0xbe, 0x4c, 0x16, 0x1d, 0x63, 0x7f, 0xf9, 0x3f, 0x00, 0x00,
	0xff, 0xff, 0xeb, 0x34, 0x76, 0xbe, 0x93, 0x04, 0x00, 0x00,
}
