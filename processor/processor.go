package processor

import (
	"bufio"
	"math/big"
	"os"
	"strconv"
	"strings"

	log "github.com/rs/zerolog"

	"energy-estimator/storage"
)

// Input defines the interface for reading processor.
type Input interface {
	Process()
}

// StdinInput implements InputSource to read processor from stdin.
type StdinInput struct {
	log   log.Logger
	store storage.Store
}

func New(l log.Logger, s storage.Store) Input {
	return &StdinInput{
		log:   l,
		store: s,
	}
}

// sanitizeInput trims spaces and leading characters such as >
func sanitizeInput(s string, lenIncompleteLine int) (string, bool, bool) {
	var empty, exitLoop bool
	// Trim leading and trailing spaces as well as leading > characters from each line
	line := strings.TrimSpace(strings.TrimLeft(s, "> "))

	// remove EOF from processing but have loop exit variable
	line, hasSuffix := strings.CutSuffix(line, "EOF")
	// if no other words in line break the loop
	if line == "" && lenIncompleteLine == 0 {
		empty = true
	}

	if hasSuffix {
		exitLoop = true
	}

	return line, exitLoop, empty

}

// checks if string is empty
func isEmpty(s string) bool {
	if len(s) == 0 {
		return true
	}
	return false
}

// Process reads processor from stdin and processes the messages
func (std *StdinInput) Process() {
	s := bufio.NewScanner(os.Stdin)

	var exitLoop, emptyString bool
	var line string
	var incompleteLine []string
	for s.Scan() {

		if s.Err() != nil {
			std.log.Error().Err(s.Err()).Msg("error parsing message")
			break
		}

		if isEmpty(s.Text()) {
			continue
		}

		line, exitLoop, emptyString = sanitizeInput(s.Text(), len(incompleteLine))
		// if input sanitized and EOF removed, no other words in the line so break the loop
		if emptyString {
			break
		}

		input := strings.Split(strings.TrimSpace(line), " ")

		// a complete line is either of the form
		// epochTs TurnOff or epochTs Delta dimmer_value

		// incomplete lines can be of the form
		// epochTs or TurnOff or Delta or dimmer_value or epochTs Delta or Delta dimmer_value
		if len(input) < 2 { // can only be 1 as we deal with empty lines above
			if input[0] != "TurnOff" && input[0] != "Delta" {
				incompleteLine = append(incompleteLine, input[0])

				continue
			}
			// alternatives are TurnOff, Delta or dimmer_value
			incompleteLine = append(incompleteLine, input[0])

			// delta still needs the dimmer value
			if input[0] == "Delta" {
				continue
			}

		}

		// when incomplete line is of the form epochTs Delta
		if len(input) == 2 {
			if input[1] == "Delta" {
				incompleteLine = append(incompleteLine, input...)
				continue
			}
		}

		if len(incompleteLine) == 2 && incompleteLine[1] == "TurnOff" {
			input = incompleteLine
			incompleteLine = nil
		}

		if len(incompleteLine) == 3 && incompleteLine[1] == "Delta" {
			input = incompleteLine
			incompleteLine = nil
		}

		// parse message fields
		action, ts, delta, err := parseMessage(input)
		if err != nil {
			std.log.Error().Err(err).Msg("error parsing message")
			continue
		}

		// store messages
		var msg *storage.Message
		switch action {
		case "turnoff":
			msg = &storage.Message{
				EpochSecTs: ts,
				Action:     action,
			}
			std.store.Put(msg)
		case "delta":
			msg = &storage.Message{
				EpochSecTs: ts,
				Action:     action,
				Delta:      delta,
			}
			std.store.Put(msg)
		default:
			continue
		}

		if exitLoop {
			break
		}
	}
}

// parseMessage converts string values from raw message to standard types
func parseMessage(input []string) (string, int64, float64, error) {
	var action string
	var ts int64
	var delta float64
	var err error

	if len(input) > 1 {
		action = strings.ToLower(input[1])
		ts, err = strconv.ParseInt(input[0], 10, 64)
		if err != nil {
			return "", 0, 0.0, err
		}
	}

	if len(input) > 2 {
		fv, err := strconv.ParseFloat(input[2], 64)
		if err != nil {
			return "", 0, 0.0, err
		}

		delta, _ = big.NewFloat(0).SetPrec(200).SetFloat64(fv).Float64()
	}

	return action, ts, delta, nil
}

// Below is an example of an extension if we were to implement pubsub
// PubSub implements InputSource to read processor from a pub/sub messaging system.
//type PubSubInput struct {
//	Topic string
//}

// ReadInput reads from a pub/sub messaging system.
// We would have liked to implement the read method in the interface above
//func (p *PubSub) Read() (*message, error) {
//	// Implementation to read from pub/sub topic
//	return nil, nil
//}

//func (p *PubSub) Process {
//	// Implementation of processor from pub/sub
//	return nil, nil
//}
