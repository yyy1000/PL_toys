package main

import "fmt"

// unsigned numbers tend to be used only when their bitwise operators or
// peculiar arithmetic operators are required.

// avoid conversions in which the operand is out of range
// for the target type, because the behavior depends on the implementation.

// very small or very large numbers are better written in scientific notation

type Currency int

const (
	USD Currency = iota
	RMB
	EUP
)



func main(){
	const GoUsage = `Go is a tool for managing Go source code.
Usage:
go command [arguments]
...`
	fmt.Print(GoUsage)
	symbol := [...]string{USD: "$", EUP: "9", RMB: "¥"}
	fmt.Println(RMB,symbol[RMB])
}
