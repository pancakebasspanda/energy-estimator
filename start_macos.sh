#!/bin/bash

# Execute Makefile to build Linux and macOS executables
make build-macos

# Run macOS executable 1
echo "./energy-estimator <<EOF
      1544206562 TurnOff
      1544206563 Delta +0.5
      1544213763 TurnOff
      EOF"
./energy-estimator <<EOF
1544206562 TurnOff
1544206563 Delta +0.5
1544213763 TurnOff
EOF

echo "\n"

# Run macOS executable 2
echo "./energy-estimator <<EOF
      > 1544206562 TurnOff
      > 1544206563 Delta +0.5
      > 1544210163 Delta -0.25
      > 1544210163 Delta -0.25
      > 1544211963 Delta +0.75
      > 1544213763
      TurnOff EOF"
./energy-estimator <<EOF
> 1544206562 TurnOff
> 1544206563 Delta +0.5
> 1544210163 Delta -0.25
> 1544210163 Delta -0.25
> 1544211963 Delta +0.75
> 1544213763
TurnOff EOF