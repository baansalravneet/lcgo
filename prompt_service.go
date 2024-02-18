package main

import (
    "strings"
    "bufio"
    "os"
    "fmt"
)

func getInput(prompt string) string {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print(prompt)
    input, _ := reader.ReadString('\n')
    return strings.TrimSpace(input)
}
