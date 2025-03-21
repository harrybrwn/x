package types

// Uint is a type constraint for unsigned integers.
type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Sint is a type constraint for signed integers.
type Sint interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Int is a type constraint for any integer type.
type Int interface {
	Sint | Uint
}

// Float is a type constraint for any floating point number type.
type Float interface {
	~float32 | ~float64
}

// Number is a type constraint for any number type.
type Number interface {
	Sint | Uint | Float
}
