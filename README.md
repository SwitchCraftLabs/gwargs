# gwargs

Command line argument parsing library for Go.

The Parse function takes in a struct pointer, parses named command line arguments and flags from
`os.Args` and populates the provided struct with the corresponding values. Parsing and type-casting
errors are reported immediately. Float and integer under and overflow are prevented and will result
in an error.

Gwargs currently supports field types of `string`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`,
`uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, and `float64`.
