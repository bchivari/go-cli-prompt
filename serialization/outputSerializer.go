package serialization

// OutputSerializer defines a function used to convert from an input string into the desired output type
type OutputSerializer func(string) (interface{}, error)
