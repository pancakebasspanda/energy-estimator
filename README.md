# energy-estimator

## About
Energy estimator is a cli tool that can estimate the power draw of a single dimmable smart light
bulb. 

When fully lit, the light consumes 5W, and this then drops linearly as the dimmer is turned down.
Internally, the light represents its dimmer value as a floating point number between 0.0 and 1.0,
inclusive.

The light outputs a message whenever someone adjusts it. Each message contains a
timestamp from the light bulb’s internal clock (in seconds since the start of 1970). There are two
types of messages. A TurnOff message indicates that the light has been turned off
completely. A Delta message indicates that the brightness has been adjusted; it includes a
value for the change in the dimmer value (a floating point number between -1.0 and +1.0
inclusive).

The command line tool reads messages from stdin
until it reaches an EOF. It should then print the estimated energy consumed in watt-hours and
exit. Here is an example session:

```shell
$ energy-estimator <<EOF
> 1544206562 TurnOff
> 1544206563 Delta +0.5
> 1544210163 TurnOff
> EOF
Estimated energy used: 2.5 Wh
```

## Assumptions
* Every stream starts and ends with a TurnOff message.
* The protocol used to transmit the messages is unreliable and messages being duplicated, lost and/or delivered out of order.
* The Estimated energy between the first TurnOff and  first delta us negligible (see explaination in the next section)

## Sample Input Logical Breakdown
```text
$ energy-estimator <<EOF
> 1544206562 TurnOff				17872 days, 18 hours, 16 minutes and 2 seconds.
> 1544206563 Delta +0.5 			17872 days, 18 hours, 16 minutes and 3 seconds.
> 1544210163 Delta -0.25			17872 days, 19 hours, 16 minutes and 3 seconds.
> 1544210163 Delta -0.25            17872 days, 19 hours, 16 minutes and 3 seconds. (Duplicate of previous event)
> 1544211963 Delta +0.75 			17872 days, 19 hours, 46 minutes and 3 seconds.
> 1544213763 						17872 days, 20 hours, 16 minutes and 3 seconds
TurnOff EOF
Estimated energy used: 5.625 Wh
```
```text
1. 1544206562 TurnOff to 1544206563 Delta +0.5
time difference : 1 second
energy cost : 5w in 1 hr
?w  in (1/60) hr
? = (5*0.5)*(1/60) =  0.001388888889 Wh (negligible)

2. 1544206563 Delta +0.5 to 1544210163 Delta -0.25
time difference : 1h
energy cost : At a dimmer value of 0.5 the effective power usage will be 5 * 0.5 = 2.5 W
2.5W * 1h = 2.5 Wh

3. 1544210163 Delta -0.25 to 1544211963 Delta +0.75
time difference : (46-16) = 30 mins => 0.5h
energy cost :  dimmer value previously at 0.5 but then adjusted by -0.25 so effective position is
0.5 - 0.25 = 0.25 which relates to a power rating of 0.25 * 5 = 1.25 W
1.25 W in 1 h
? in 0.5 h
? = 1.25*0.5 => 0.625 Wh

4. 1544211963 Delta +0.75 to 1544213763 TurnOff EOF
time difference : 20h 16m - 19h 46m = 30 mins => 0.5 h
energy cost : dimmer value previously at 0.25 but then adjusted by +0.75 so effective position is 1.0
which relates to a power rating of 1 * 5 = 5 W
5 W in 1 h
? in 0.5 h
? = 5 * 0.5 = 2.5 Wh

Est. total energy used = 2.5 + 0.625 + 2.5 = 5.625 Wh

```


## Local Machine Run using script start.sh
### macOS

```shell
$ sh start_macos.sh                                                           
Building macOS executable
env GOOS=darwin GOARCH=amd64 go build -o energy-estimator ./cmd/energy-estimator
./energy-estimator <<EOF
      1544206562 TurnOff
      1544206563 Delta +0.5
      1544213763 TurnOff
      EOF
Estimated energy used: 5 Wh

./energy-estimator <<EOF
      > 1544206562 TurnOff
      > 1544206563 Delta +0.5
      > 1544210163 Delta -0.25
      > 1544210163 Delta -0.25
      > 1544211963 Delta +0.75
      > 1544213763
      TurnOff EOF
Estimated energy used: 5.625 Wh                                                                                                                     ➜  energy-estimator git:(master) ✗ 

```

## Running tests
### Unit Tests
First depending on you operating system run either 
```shell
make build-linux
```

or 

```shell
make build-macos
```

then 

```shell
➜  energy-estimator ✗ make unit-test                                                              
Running unit tests
go test -v energy-estimator/cmd/energy-estimator energy-estimator/processor energy-estimator/storage -race
?       energy-estimator/cmd/energy-estimator   [no test files]
=== RUN   Test_sanitizeInput
=== RUN   Test_sanitizeInput/successfully_sanitizes_input_when_line_is_complete
=== RUN   Test_sanitizeInput/successfully_sanitizes_input_when_line_is_incomplete
=== RUN   Test_sanitizeInput/successfully_sanitizes_input_when_EOF_included_
--- PASS: Test_sanitizeInput (0.00s)
    --- PASS: Test_sanitizeInput/successfully_sanitizes_input_when_line_is_complete (0.00s)
    --- PASS: Test_sanitizeInput/successfully_sanitizes_input_when_line_is_incomplete (0.00s)
    --- PASS: Test_sanitizeInput/successfully_sanitizes_input_when_EOF_included_ (0.00s)
=== RUN   Test_isEmpty
--- PASS: Test_isEmpty (0.00s)
PASS
ok      energy-estimator/processor      (cached)
=== RUN   TestSmartLightStore_CalculateEstimatedPowerConsumption
=== RUN   TestSmartLightStore_CalculateEstimatedPowerConsumption/expect_estimate_of_2.5
=== RUN   TestSmartLightStore_CalculateEstimatedPowerConsumption/sum_of_delta_values_less_than_dimmer_range
=== RUN   TestSmartLightStore_CalculateEstimatedPowerConsumption/sum_of_delta_values_larger_than_dimmer_range
=== RUN   TestSmartLightStore_CalculateEstimatedPowerConsumption/expect_estimate_of_5.625
=== RUN   TestSmartLightStore_CalculateEstimatedPowerConsumption/expect_estimate_of_1
--- PASS: TestSmartLightStore_CalculateEstimatedPowerConsumption (0.00s)
    --- PASS: TestSmartLightStore_CalculateEstimatedPowerConsumption/expect_estimate_of_2.5 (0.00s)
    --- PASS: TestSmartLightStore_CalculateEstimatedPowerConsumption/sum_of_delta_values_less_than_dimmer_range (0.00s)
    --- PASS: TestSmartLightStore_CalculateEstimatedPowerConsumption/sum_of_delta_values_larger_than_dimmer_range (0.00s)
    --- PASS: TestSmartLightStore_CalculateEstimatedPowerConsumption/expect_estimate_of_5.625 (0.00s)
    --- PASS: TestSmartLightStore_CalculateEstimatedPowerConsumption/expect_estimate_of_1 (0.00s)
=== RUN   TestSmartLightStore_Put_Get
=== RUN   TestSmartLightStore_Put_Get/successfully_puts_/_gets_TurnOff_message_into_local_storage_and_is_ordered_by_timestamp_asc
=== RUN   TestSmartLightStore_Put_Get/successfully_puts_/_gets_Delta_message_into_local_storage_and_is_ordered_by_timestamp_asc
--- PASS: TestSmartLightStore_Put_Get (0.00s)
    --- PASS: TestSmartLightStore_Put_Get/successfully_puts_/_gets_TurnOff_message_into_local_storage_and_is_ordered_by_timestamp_asc (0.00s)
    --- PASS: TestSmartLightStore_Put_Get/successfully_puts_/_gets_Delta_message_into_local_storage_and_is_ordered_by_timestamp_asc (0.00s)
PASS
ok      energy-estimator/storage        (cached)
➜  energy-estimator ✗ 
 
```
### Integration Tests
```shell
➜  energy-estimator ✗ make integration-test 
Running integration tests
go test ./integration-test -v
=== RUN   TestCLIIntegration
=== RUN   TestCLIIntegration/Estimated_energy_of_2.5_Wh
=== RUN   TestCLIIntegration/Estimated_energy_of_5.625_Wh
=== RUN   TestCLIIntegration/Input_/_messages_are_multi_line
=== RUN   TestCLIIntegration/Input_/_messages_are_multi_line#01
--- PASS: TestCLIIntegration (0.03s)
    --- PASS: TestCLIIntegration/Estimated_energy_of_2.5_Wh (0.01s)
    --- PASS: TestCLIIntegration/Estimated_energy_of_5.625_Wh (0.01s)
    --- PASS: TestCLIIntegration/Input_/_messages_are_multi_line (0.01s)
    --- PASS: TestCLIIntegration/Input_/_messages_are_multi_line#01 (0.01s)
PASS
ok      energy-estimator/integration-test       (cached)
➜  energy-estimator ✗ 

```
### All Tests
```shell
➜  energy-estimator ✗ make all-tests       

```

## Packages
`/integration_test`: intergration tests run by using the os.exec go packages functionality

`/cmd`: main.go + setup directly related to logging and ingestion of stdin messages

`/storage`: in app memory storage solution for deduping and keeping track of the messages

`/processor`: core logic around message processing and input parsing