package storage

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

var _ Store = (*SmartLightStore)(nil)

type Store interface {
	Put(msg *Message)
	Get(ts int64) (*Message, error)
	PrintStore()
	CalculateEstimatedPowerConsumption() float64
}

// SmartLightStore is an app memory store for smart light messages
type SmartLightStore struct {
	MsgStore map[int64]*Message // quick reads when accessing data
	TsStore  []int64            // track for duplicates and order
}

func New() Store {
	return &SmartLightStore{
		MsgStore: make(map[int64]*Message, 0),
		TsStore:  make([]int64, 0),
	}
}

// Put inserts data into our in memory storage store
func (s *SmartLightStore) Put(msg *Message) {
	// search function to find index where the new element will be inserted
	index := sort.Search(len(s.TsStore), func(i int) bool {
		return s.TsStore[i] >= msg.EpochSecTs
	})

	// If the element already exists, return the slice as it is
	if index < len(s.TsStore) && s.TsStore[index] == msg.EpochSecTs {
		return
	}

	// Allocate space for the new element
	s.TsStore = append(s.TsStore, 0)

	// Shift elements to the right
	copy(s.TsStore[index+1:], s.TsStore[index:])

	// insert into the correct index (increasing order)
	s.TsStore[index] = msg.EpochSecTs

	// create entry in map
	s.MsgStore[msg.EpochSecTs] = msg

	return
}

// Get retrieves data from our in memory storage store
func (s *SmartLightStore) Get(ts int64) (*Message, error) {
	if msg, ok := s.MsgStore[ts]; ok {
		return msg, nil
	}

	return nil, errors.New("smart light storage does not exist")
}

// PrintStore is a helper method to printStore if needed
func (s *SmartLightStore) PrintStore() {
	for _, ts := range s.TsStore {
		fmt.Println(s.MsgStore[ts])
	}
}

// CalculateEstimatedPowerConsumption calculates the total power consumption
// by evaluating the dimmer value and time at the specified power
func (s *SmartLightStore) CalculateEstimatedPowerConsumption() float64 {
	energyTotal := 0.0
	prevDimmer := 0.0
	var dimmer float64
	for i := 1; i < len(s.TsStore)-1; i++ {
		// ignore calculation for the period from the first turn to first delta as 0 power will be used
		t1 := time.Unix(s.MsgStore[s.TsStore[i]].EpochSecTs, 0)
		t2 := time.Unix(s.MsgStore[s.TsStore[i+1]].EpochSecTs, 0)

		tsDiff := t2.Sub(t1).Hours()

		// the lights dimmer value is a floating point between 0.0 and 1.0
		if prevDimmer+s.MsgStore[s.TsStore[i]].Delta <= 0.0 {
			dimmer = 0.0
		} else if prevDimmer+s.MsgStore[s.TsStore[i]].Delta >= 1.0 {
			dimmer = 1.0
		} else {
			dimmer = prevDimmer + s.MsgStore[s.TsStore[i]].Delta
		}

		prevDimmer = dimmer
		energyTotal += dimmer * 5 * tsDiff
	}

	return FixedPrecision(energyTotal, 3)
}

// FixedPrecision returns the estimated energy with simple precision formula below. Will fail if numbers flow over float64
func FixedPrecision(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Round(num*output) / output
}
