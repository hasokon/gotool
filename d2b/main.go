package main

import (
	"flag"
	"fmt"
	"strconv"
)

func decimal2binary(inputvalue string) {
	value, err := strconv.Atoi(inputvalue)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer fmt.Println("")
	if value == 0 {
		defer fmt.Print("0")
		return
	}

	for value > 0 {
		defer fmt.Print(value % 2)
		value = value / 2
	}
}

func binary2decimal(inputvalue string) {
	value, err := strconv.ParseInt(inputvalue, 2, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)
}

func main() {
	var (
		isReverse bool
	)
	
	flag.BoolVar(&isReverse, "r", false, "convert binary to decimal")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage : d2b [OPTION] <Decimal value / Binary value>")
		return
	}

	switch (true) {
		case isReverse :
			binary2decimal(flag.Arg(0))
		default :
			decimal2binary(flag.Arg(0))
	}
}
