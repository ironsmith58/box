package main

import (
	"fmt"

	"github.com/pborman/getopt/v2"
)

type BoxType struct {
	topleft, topright, botleft, botright, leftside, rightside, topside, botside string
}

var Boxes = map[string]BoxType{
	"single": {"┌", "┐", "└", "┘", "│", "│", "─", "─"},
	"double": {"╔", "╗", "╚", "╝", "║", "║", "═", "═"},
	"round":  {"╭", "╮", "╰", "╯", "│", "│", "─", "─"},
	"bold":   {"┏", "┓", "┗", "┛", "┃", "┃", "━", "━"},
	"shadow": {"┌", "┐", "└", "┘", "│", "│", "─", "─"},
	"simple": {".-", "-.", "`-", "-`", "|", "|", "-", "-"},
}

var (
	verbose = getopt.BoolLong("verbose", 'v', "Print more information")
)

func main() {
	getopt.Parse()

	if verbose != nil && *verbose {
		fmt.Println("Hello, World!")
	}
}
