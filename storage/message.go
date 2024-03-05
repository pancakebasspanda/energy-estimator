package storage

type Message struct {
	EpochSecTs int64
	Action     string
	Delta      float64
}
