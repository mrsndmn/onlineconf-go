package expr

import "strings"

type (
	// A Kind defines the conceptual type that a DataType represents.
	Kind uint

	// DataType is the common interface to all types.
	DataType interface {
		// Kind of data type, one of the Kind enum.
		Kind() Kind
		// Name returns the type name.
		Name() string
	}

	// Primitive is the type for null, boolean, integer, number, string, and time.
	Primitive Kind

	// JSON is interface for JSON params
	JSON interface {
		UnmarshalJSON([]byte)
	}

	OnlineConfPath string
)

const (
	// BooleanKind represents a boolean.
	BooleanKind Kind = iota + 1
	// IntKind represents a signed integer.
	IntKind
	// Int32Kind represents a signed 32-bit integer.
	Int32Kind
	// Int64Kind represents a signed 64-bit integer.
	Int64Kind
	// UIntKind represents an unsigned integer.
	UIntKind
	// UInt32Kind represents an unsigned 32-bit integer.
	UInt32Kind
	// UInt64Kind represents an unsigned 64-bit integer.
	UInt64Kind
	// Float32Kind represents a 32-bit floating number.
	Float32Kind
	// Float64Kind represents a 64-bit floating number.
	Float64Kind
	// StringKind represents a JSON string.
	StringKind
	// BytesKind represent a series of bytes (binary data).
	BytesKind
	// JSONKind represents a JSON array.
	JSONKind
	// AnyKind represents an unknown type.
	AnyKind
)

const (
	// Boolean is the type for a JSON boolean.
	Boolean = Primitive(BooleanKind)

	// Int is the type for a signed integer.
	Int = Primitive(IntKind)

	// Int32 is the type for a signed 32-bit integer.
	Int32 = Primitive(Int32Kind)

	// Int64 is the type for a signed 64-bit integer.
	Int64 = Primitive(Int64Kind)

	// UInt is the type for an unsigned integer.
	UInt = Primitive(UIntKind)

	// UInt32 is the type for an unsigned 32-bit integer.
	UInt32 = Primitive(UInt32Kind)

	// UInt64 is the type for an unsigned 64-bit integer.
	UInt64 = Primitive(UInt64Kind)

	// Float32 is the type for a 32-bit floating number.
	Float32 = Primitive(Float32Kind)

	// Float64 is the type for a 64-bit floating number.
	Float64 = Primitive(Float64Kind)

	// String is the type for a JSON string.
	String = Primitive(StringKind)

	// Bytes is the type for binary data.
	Bytes = Primitive(BytesKind)

	// Any is the type for an arbitrary JSON value (interface{} in Go).
	Any = Primitive(AnyKind)
)

// Kind implements DataKind.
func (p Primitive) Kind() Kind { return Kind(p) }

// Name returns the type name appropriate for logging.
func (p Primitive) Name() string {
	switch p {
	case Boolean:
		return "boolean"
	case Int:
		return "int"
	case Int32:
		return "int32"
	case Int64:
		return "int64"
	case UInt:
		return "uint"
	case UInt32:
		return "uint32"
	case UInt64:
		return "uint64"
	case Float32:
		return "float32"
	case Float64:
		return "float64"
	case String:
		return "string"
	case Bytes:
		return "bytes"
	case Any:
		return "any"
	default:
		panic("unknown primitive type") // bug
	}
}

func (j JSON) Kind() Kind {
	return JSONKind
}

func (j JSON) Name() string {
	return "json"
}

func (p OnlineConfPath) IsValid() bool {
	pstr := string(p)

	return len(pstr) > 0 && (pstr == "/" || ( strings.HasPrefix(pstr, "/") &&! strings.HasSuffix(pstr, "/") ))
}