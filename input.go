package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadInputYesNo() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	fields := strings.Fields(input) // hacky way to remove whitespace characters
	if len(fields) == 0 {
		return false, fmt.Errorf("invalid input")
	}
	return strings.ToLower(fields[0]) == "y", nil
}

func ReadInputNumbers() ([]int, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	check(err)

	inputFields := strings.Fields(input)
	var subsToRemove []int
	for _, n := range inputFields {
		i, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return nil, err
		}
		subsToRemove = append(subsToRemove, int(i))
	}

	return subsToRemove, nil
}
