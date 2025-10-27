package gonv

// Compatibility wrapper functions for older test names / API.
// These delegate to the new generic functions.

// ToStringE is a compatibility wrapper for StringE[string]
func ToStringE(o any) (string, error) { return StringE[string](o) }

// ToIntE is a compatibility wrapper for IntE[int]
func ToIntE(o any) (int, error) { return IntE[int](o) }

// ToInt64E is a compatibility wrapper for IntE[int64]
func ToInt64E(o any) (int64, error) { return IntE[int64](o) }

// ToUintE is a compatibility wrapper for UintE[uint]
func ToUintE(o any) (uint, error) { return UintE[uint](o) }

// ToUint64E is a compatibility wrapper for UintE[uint64]
func ToUint64E(o any) (uint64, error) { return UintE[uint64](o) }

// ToUint32E is a compatibility wrapper for UintE[uint32]
func ToUint32E(o any) (uint32, error) { return UintE[uint32](o) }

// ToFloat64E is a compatibility wrapper for FloatE[float64]
func ToFloat64E(o any) (float64, error) { return FloatE[float64](o) }

// ToFloat32E is a compatibility wrapper for FloatE[float32]
func ToFloat32E(o any) (float32, error) { return FloatE[float32](o) }
