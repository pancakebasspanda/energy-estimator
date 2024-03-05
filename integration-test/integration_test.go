package integration_test_test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"strings"
	"testing"
)

func TestCLIIntegration(t *testing.T) {
	// Run CLI command
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "Estimated energy of 2.5 Wh",
			input: "> 1544206562 TurnOff\n" +
				"> 1544206563 Delta +0.5\n" +
				"> 1544210163 TurnOff\n" +
				"EOF",
			want: "Estimated energy used: 2.5 Wh",
		}, {
			name: "Estimated energy of 5.625 Wh",
			input: "> 1544206562 TurnOff\n" +
				"> 1544206563 Delta +0.5\n" +
				"> 1544210163 Delta -0.25\n" +
				"> 1544210163 Delta -0.25\n" +
				"> 1544211963 Delta +0.75\n" +
				"> 1544213763 TurnOff\n" +
				"EOF",
			want: "Estimated energy used: 5.625 Wh",
		},
		{
			name: "Input / messages are multi line",
			input: "> 1544206562 TurnOff\n" +
				"> 1544206563 Delta +0.5\n" +
				"> 1544210163 Delta -0.25\n" +
				"> 1544210163 Delta -0.25\n" +
				"> 1544211963 Delta +0.75\n" +
				"> 1544213763 TurnOff " +
				"EOF",
			want: "Estimated energy used: 5.625 Wh",
		},
		{
			name: "Input / messages are multi line",
			input: "> 1544206562 TurnOff\n" +
				"> 1544206563 Delta +0.5\n" +
				"> 1544210163 Delta -0.25\n" +
				"> 1544210163 Delta -0.25\n" +
				"> 1544211963 Delta +0.75\n" +
				"> 1544213763\n" +
				"TurnOff EOF",
			want: "Estimated energy used: 5.625 Wh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cmd := exec.Command("../energy-estimator", "<<EOF")

			cmd.Stdin = strings.NewReader(tt.input)

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("CLI command failed: %v", err)
			}

			assert.Equal(t, tt.want, string(output))
		})
	}
}
