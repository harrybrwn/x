package ihash

// Fn is the type definition for hash functions in the package.
type Fn func([]byte) uint64

// This hash function comes from Brian Kernighan and Dennis Ritchie's book "The
// C Programming Language". It is a simple hash function using a strange set of
// possible seeds which all constitute a pattern of 31....31...31 etc, it seems
// to be very similar to the DJB hash function.
func BKDR(data []byte) (hash uint64) {
	const seed = 131
	for _, c := range data {
		hash = (hash * seed) + uint64(c)
	}
	return hash
}

// Similar to the PJW Hash function, but tweaked for 32-bit processors. It is a
// widely used hash function on UNIX based systems.
func Elf(data []byte) (hash uint64) {
	var x uint64
	for _, c := range data {
		hash = (hash << 4) + uint64(c)
		if x = hash & 0xf0000000; x != 0 {
			hash ^= (x >> 24)
		}
		hash &= (^x)
	}
	return hash
}

func Js(data []byte) uint64 {
	hash := uint64(1315423911)
	for _, c := range data {
		hash ^= ((hash << 5) + uint64(c) + (hash >> 2))
	}
	return hash
}

func Rs(data []byte) (hash uint64) {
	var a, b uint64 = 63689, 378551
	for _, c := range data {
		hash = hash*a + uint64(c)
		a *= b
	}
	return hash
}

func Dek(data []byte) (hash uint64) {
	for _, c := range data {
		hash = ((hash << 5) ^ (hash >> 27)) ^ uint64(c)
	}
	return hash
}

func Djb2(data []byte) uint64 {
	hash := uint64(5381)
	for _, b := range data {
		hash = ((hash << 5) + hash) + uint64(b)
	}
	return hash
}

func Djbx33a(data []byte) uint64 {
	hash := uint64(5381)
	for _, b := range data {
		hash = hash*33 + uint64(b)
	}
	return hash
}

const (
	fnvOffsetBasis = 0xcbf29ce484222325
	fnvPrime       = 0x100000001b3
)

func Fnv1(data []byte) uint64 {
	hash := uint64(fnvOffsetBasis)
	for _, b := range data {
		hash *= fnvPrime
		hash ^= uint64(b)
	}
	return hash
}

func Fnv1a(data []byte) uint64 {
	hash := uint64(fnvOffsetBasis)
	for _, b := range data {
		hash ^= uint64(b)
		hash *= fnvPrime
	}
	return hash
}

func Sdbm(data []byte) (hash uint64) {
	for _, b := range data {
		hash = uint64(b) + (hash << 6) + (hash << 16) - hash
	}
	return hash
}
