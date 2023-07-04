package gorundut

import "runtime/metrics"

const (
	goMemoryClassesHeapObjectsBytes = "/memory/classes/heap/objects:bytes"
)

func GetHeapAllocTotal() MemoryBytes {
	samples := []metrics.Sample{{
		Name: goMemoryClassesHeapObjectsBytes,
	}}
	metrics.Read(samples)
	return MemoryBytes(samples[0].Value.Uint64())
}
