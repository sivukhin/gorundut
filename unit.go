package gorundut

type MemoryBytes uint64

const (
	Byte     MemoryBytes = 1
	Kilobyte             = 1024 * Byte
	Megabyte             = 1024 * Kilobyte
	Gigabyte             = 1024 * Megabyte
	Terabyte             = 1024 * Gigabyte
)
