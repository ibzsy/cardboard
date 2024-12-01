package main

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/ibzsy/cardboard/object"
	"github.com/ibzsy/cardboard/repl"
)

func main() {
	var buffer bytes.Buffer
	var env = object.CreateEnvironment()

	for {
		repl.StartREPL(&buffer, env)
		// Create a Scanner to read each string (line)
		scanner := bufio.NewScanner(bytes.NewReader(buffer.Bytes()))

		// Iterate over each line in the buffer
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // Print each string
		}

		// Check for any errors during scanning
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading buffer:", err)
		}
		buffer.Reset()
	}
}
