package calculator

const NO_CLEANUP = -1
const ROUTINE_CLEANUP = -2

type Calculator interface {
	ShouldCollectGarbageOnCreate() bool
	CalculateStoreSize() int64
	CalculateGCThreshold() int64
}

type modernCalculator struct {
	reservedSpace int64
	diskSize      int64
	minStoreSize  int64
}

func NewModernCalculator(reservedSpace, diskSize, minStoreSize int64) Calculator {
	return &modernCalculator{
		reservedSpace: reservedSpace,
		diskSize:      diskSize,
		minStoreSize:  minStoreSize,
	}
}

func (m modernCalculator) ShouldCollectGarbageOnCreate() bool {
	return m.reservedSpace != NO_CLEANUP
}

func (m modernCalculator) CalculateStoreSize() int64 {
	storeSize := m.diskSize - positiveOrZero(m.reservedSpace)
	if storeSize > m.minStoreSize {
		return storeSize
	}

	return m.diskSize
}

func (m modernCalculator) CalculateGCThreshold() int64 {
	if m.reservedSpace == ROUTINE_CLEANUP {
		return 0
	}

	return positiveOrZero(m.diskSize - positiveOrZero(m.reservedSpace))
}

type oldFashionedCalculator struct {
	diskSize    int64
	gcThreshold int64
}

func NewOldFashionedCalculator(diskSize, gardenGcThreshold, grootfsGcThreshold int64) *oldFashionedCalculator {
	effectiveGcThreshold := grootfsGcThreshold
	if gardenGcThreshold > -1 {
		effectiveGcThreshold = gardenGcThreshold
	}

	return &oldFashionedCalculator{
		diskSize:    diskSize,
		gcThreshold: effectiveGcThreshold,
	}
}

func (o oldFashionedCalculator) ShouldCollectGarbageOnCreate() bool {
	return o.gcThreshold >= 0
}

func (o oldFashionedCalculator) CalculateStoreSize() int64 {
	return o.diskSize
}

func (o oldFashionedCalculator) CalculateGCThreshold() int64 {
	if o.gcThreshold > 0 {
		return o.gcThreshold
	}
	return o.diskSize + 1
}

func positiveOrZero(n int64) int64 {
	if n < 0 {
		return 0
	}

	return n
}
