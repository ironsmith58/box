package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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

var verbose = getopt.BoolLong("verbose", 'v', "Print more information")
var listBoxTypes = getopt.BoolLong("list", 'l', "List all available box types")
var boxTypeName = getopt.StringLong("box", 'b', "single", "Select box type by name")

func wrapInBox(lines []string, boxType BoxType) []string {
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}
	var boxed []string
	// Top border
	boxed = append(boxed, boxType.topleft+strings.Repeat(boxType.topside, maxLen+2)+boxType.topright)
	// Content lines
	for _, line := range lines {
		padding := maxLen - len(line)
		boxed = append(boxed, boxType.leftside+" "+line+strings.Repeat(" ", padding)+" "+boxType.rightside)
	}
	// Bottom border
	boxed = append(boxed, boxType.botleft+strings.Repeat(boxType.botside, maxLen+2)+boxType.botright)
	return boxed
}

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
		fmt.Fprintf(os.Stderr, "Unknown box type: %s, using 'single'\n", *boxTypeName)
	}

	boxed := wrapInBox(lines, boxType)
	for _, line := range boxed {
		fmt.Println(line)
	}
}
