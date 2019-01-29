package main

type Config struct {
	Checks []ICheck `json:"checks"`
}

type ICheck interface {
	Validate() bool
}
