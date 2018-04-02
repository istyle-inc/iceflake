package app

// IDGenerator interface of generator generates each Unique ID
type IDGenerator interface {
	Generate() (uint64, error)
}

type IceFlakeGenerator struct {
	w int
}

func NewIDGenerator(workerID int) IDGenerator {
	return &IceFlakeGenerator{
		w: workerID,
	}
}

func (g IceFlakeGenerator) Generate() (uint64, error) {
	return 0, nil
}
