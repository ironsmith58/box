package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pborman/getopt/v2"
)

// BoxType defines the characters used for each part of a box.
type BoxType struct {
	topleft, topright, botleft, botright, leftside, rightside, topside, botside string
}

var Boxes = map[string]BoxType{
	"single":      {"┌", "┐", "└", "┘", "│", "│", "─", "─"},
	"double":      {"╔", "╗", "╚", "╝", "║", "║", "═", "═"},
	"round":       {"╭", "╮", "╰", "╯", "│", "│", "─", "─"},
	"bold":        {"┏", "┓", "┗", "┛", "┃", "┃", "━", "━"},
	"shadow":      {"┌", "┐", "└", "┘", "│", "│", "─", "─"},
	"simple":      {".-", "-.", "`-", "-`", "|", "|", "-", "-"},
	"triple":      {"╓", "╖", "╙", "╜", "║", "║", "═", "═"}, // Triple line corners
	"block":       {"█", "█", "█", "█", "█", "█", "█", "█"}, // Full block
	"dotted":      {"┈", "┈", "┈", "┈", "┊", "┊", "┈", "┈"}, // Dotted lines
	"dash":        {"┄", "┄", "┄", "┄", "┆", "┆", "┄", "┄"}, // Dashed lines
}

var verbose = getopt.BoolLong("verbose", 'v', "Print more information")
var listBoxTypes = getopt.BoolLong("list", 'l', "List all available box types")
var boxTypeName = getopt.StringLong("box", 'b', "single", "Select box type by name")
var indentBy = getopt.IntLong("indent", 'i', 0, "Indent box by N spaces")

// wrapInBox returns a slice of strings representing the input lines wrapped in a box.
// It calculates the maximum line length, expands each line to fit, and adds borders.
// Handles multi-character corners and sides for proper alignment.
func wrapInBox(lines []string, boxType BoxType) []string {
	// Determine the width of each border character
	leftCornerLen := len([]rune(boxType.topleft))
	rightCornerLen := len([]rune(boxType.topright))
	leftSideLen := len([]rune(boxType.leftside))
	rightSideLen := len([]rune(boxType.rightside))

	// Calculate the maximum line length
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	var boxed []string
	// Top border
	top := boxType.topleft +
		strings.Repeat(boxType.topside, maxLen+leftSideLen+rightSideLen) +
		boxType.topright
	boxed = append(boxed, top)

	// Content lines
	for _, line := range lines {
		padding := maxLen - len(line)
		content := boxType.leftside + strings.Repeat(" ", leftCornerLen) + line + strings.Repeat(" ", padding) + strings.Repeat(" ", rightCornerLen) + boxType.rightside
		boxed = append(boxed, content)
	}

	// Bottom border
	bottom := boxType.botleft +
		strings.Repeat(boxType.botside, maxLen+leftSideLen+rightSideLen) +
		boxType.botright
	boxed = append(boxed, bottom)

	return boxed
}

// expandTabs replaces tab characters in a string with spaces up to the next tab stop.
// This ensures consistent alignment for tabbed text.
func expandTabs(line string, tabstop int) string {
	var result strings.Builder
	col := 0
	for _, r := range line {
		if r == '\t' {
			spaces := tabstop - (col % tabstop)
			result.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			result.WriteRune(r)
			// For simplicity, treat all runes as width 1
			col++
		}
	}
	return result.String()
}

// main parses command-line options, reads input, wraps it in a box, and prints the result.
// It supports selecting box types, listing available types, and verbose output.
func main() {
	getopt.Parse()

	if listBoxTypes != nil && *listBoxTypes {
		fmt.Println("Available box types:")
		for name := range Boxes {
			fmt.Println(name)
		}
		return
	}

	args := getopt.Args()
	if verbose != nil && *verbose {
		fmt.Println("Hello, World!")
	}

	var lines []string
	if len(args) > 0 {
		input := strings.Join(args, " ")
		if unquoted, err := strconv.Unquote(`"` + input + `"`); err == nil {
			for _, l := range strings.Split(unquoted, "\n") {
				lines = append(lines, expandTabs(l, 8))
			}
		} else {
			lines = []string{expandTabs(input, 8)}
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			if unquoted, err := strconv.Unquote(`"` + text + `"`); err == nil {
				lines = append(lines, expandTabs(unquoted, 8))
			} else {
				lines = append(lines, expandTabs(text, 8))
			}
		}
	}

	boxType := Boxes["single"]
	if t, ok := Boxes[*boxTypeName]; ok {
		boxType = t
	} else if *boxTypeName != "single" {
		fmt.Fprintf(os.Stderr, "Unknown box type: %s\n", *boxTypeName)
		os.Exit(1)
	}

	indent := strings.Repeat(" ", *indentBy)
	boxed := wrapInBox(lines, boxType)
	for _, line := range boxed {
		fmt.Println(indent,line)
	}
}
