package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmartLightStore_CalculateEstimatedPowerConsumption(t *testing.T) {
	tests := []struct {
		name     string
		MsgStore map[int64]*Message
		TsStore  []int64
		want     float64
	}{
		{
			name: "expect estimate of 2.5",
			MsgStore: map[int64]*Message{
				1544206562: {
					EpochSecTs: 1544206562,
					Action:     "TurnOff",
				},
				1544206563: {
					EpochSecTs: 1544206563,
					Action:     "Delta",
					Delta:      +0.5,
				},
				1544210163: {
					EpochSecTs: 1544210163,
					Action:     "TurnOff",
				},
			},
			TsStore: []int64{1544206562, 1544206563, 1544210163},
			want:    2.5,
		},
		{
			// here we increment delta from 0 to 0.5
			// then decrease delta by 0.6, but we can't go below the limit of 0 for the bulb,
			// so it is essentially the same as a switch off
			name: "sum of delta values less than dimmer range",
			MsgStore: map[int64]*Message{
				1544206562: {
					EpochSecTs: 1544206562,
					Action:     "TurnOff",
				},
				1544206563: {
					EpochSecTs: 1544206563,
					Action:     "Delta",
					Delta:      +0.5,
				},
				1544210163: {
					EpochSecTs: 1544210163,
					Action:     "Delta",
					Delta:      -0.6,
				},
				1544213763: {
					EpochSecTs: 1544213763,
					Action:     "TurnOff",
				},
			},
			TsStore: []int64{1544206562, 1544206563, 1544210163, 1544213763},
			want:    2.5,
		},
		{
			// here we increment delta from 0 to 0.5
			// then increase delta by 0.6, but we can't go above the limit of 1 for the bulb,
			// so it is essentially the same as a switch off
			name: "sum of delta values larger than dimmer range",
			MsgStore: map[int64]*Message{
				1544206562: {
					EpochSecTs: 1544206562,
					Action:     "TurnOff",
				},
				1544206563: {
					EpochSecTs: 1544206563,
					Action:     "Delta",
					Delta:      +0.5,
				},
				1544210163: {
					EpochSecTs: 1544210163,
					Action:     "Delta",
					Delta:      +0.6,
				},
				1544213763: {
					EpochSecTs: 1544213763,
					Action:     "TurnOff",
				},
			},
			TsStore: []int64{1544206562, 1544206563, 1544210163, 1544213763},
			want:    7.5,
		},
		{
			name: "expect estimate of 5.625",
			MsgStore: map[int64]*Message{
				1544206562: {
					EpochSecTs: 1544206562,
					Action:     "TurnOff",
				},
				1544206563: {
					EpochSecTs: 1544206563,
					Action:     "Delta",
					Delta:      0.5,
				},
				1544210163: {
					EpochSecTs: 1544210163,
					Action:     "Delta",
					Delta:      -0.25,
				},
				1544211963: {
					EpochSecTs: 1544211963,
					Action:     "Delta",
					Delta:      0.75,
				},
				1544213763: {
					EpochSecTs: 1544213763,
					Action:     "TurnOff",
				},
			},

			TsStore: []int64{1544206562, 1544206563, 1544210163, 1544211963, 1544213763},
			want:    5.625,
		},
		{
			name: "expect estimate of 1",
			MsgStore: map[int64]*Message{
				1544206700: { // 17872 days, 18 hours, 18 minutes and 20 seconds.
					EpochSecTs: 1544206700,
					Action:     "TurnOff",
				},
				1544216700: { // 17872 days, 21 hours, 5 minutes and 0 seconds.
					EpochSecTs: 1544216700,
					Action:     "Delta",
					Delta:      0.3,
				},
				1544226700: { // 17872 days, 23 hours, 51 minutes and 40 seconds.
					EpochSecTs: 1544226700,
					Action:     "Delta",
					Delta:      -0.2,
				},
				1544229700: { // 17873 days, 0 hours, 41 minutes and 40 seconds.
					EpochSecTs: 1544229700,
					Action:     "Delta",
					Delta:      0.1,
				},
				1544236700: { // 17873 days, 2 hours, 38 minutes and 20 seconds
					EpochSecTs: 1544236700,
					Action:     "TurnOff",
				},
			},

			TsStore: []int64{1544206700, 1544216700, 1544226700, 1544229700, 1544236700},
			want:    6.528,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SmartLightStore{
				MsgStore: tt.MsgStore,
				TsStore:  tt.TsStore,
			}
			assert.Equal(t, tt.want, s.CalculateEstimatedPowerConsumption())
		})
	}
}

func TestSmartLightStore_Put_Get(t *testing.T) {
	tests := []struct {
		name     string
		MsgStore map[int64]*Message
		TsStore  []int64
		msg      *Message
		err      error
	}{
		{
			name: "successfully puts / gets TurnOff message into local storage and is ordered by timestamp asc",
			MsgStore: map[int64]*Message{
				1544206563: {
					EpochSecTs: 1544206563,
					Action:     "Delta",
					Delta:      0.3,
				},
			},
			TsStore: []int64{1544206563},
			msg: &Message{
				EpochSecTs: 1544206562,
				Action:     "TurnOff",
			},
		},
		{
			name: "successfully puts / gets Delta message into local storage and is ordered by timestamp asc",
			MsgStore: map[int64]*Message{
				1544210163: {
					EpochSecTs: 1544210163,
					Action:     "Delta",
					Delta:      0.3,
				},
			},
			TsStore: []int64{1544210163},
			msg: &Message{
				EpochSecTs: 1544206562,
				Action:     "Delta",
				Delta:      -0.6,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SmartLightStore{
				MsgStore: tt.MsgStore,
				TsStore:  tt.TsStore,
			}
			s.Put(tt.msg)
			// check that the msg was inserted
			insertedMsg, err := s.Get(tt.msg.EpochSecTs)
			assert.NoError(t, err)
			assert.Equal(t, tt.msg, insertedMsg)
			// check the order in TsStore (should be inserted at the start)
			assert.Equal(t, tt.msg.EpochSecTs, s.TsStore[0])
		})
	}
}
