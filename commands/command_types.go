package gitmedia

import (
	_ "../gitconfig"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TypesCommand struct {
	*Command
}

func (c *TypesCommand) Run() {
	var sub string
	if len(c.SubCommands) > 0 {
		sub = c.SubCommands[0]
	}

	switch sub {
	case "add":
		c.addType()
	case "remove":
		fmt.Println("Removing type")
	default:
		c.listTypes()
	}

}

func (c *TypesCommand) addType() {
	if len(c.SubCommands) < 2 {
		fmt.Println("git media types add <type> [type]*")
		return
	}

	knownTypes := findTypes()
	attributesFile, err := os.OpenFile(".gitattributes", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println("Error opening .gitattributes file")
		return
	}

	for _, t := range c.SubCommands[1:] {
		isKnownType := false
		for _, k := range knownTypes {
			if t == k {
				isKnownType = true
			}
		}

		if isKnownType {
			fmt.Println(t, "already supported")
			continue
		}

		_, err := attributesFile.WriteString(fmt.Sprintf("%s filter=media -crlf", t))
		if err != nil {
			fmt.Println("Error adding type", t)
			continue
		}
		fmt.Println("Adding type", t)
	}

	attributesFile.Close()
}

func (c *TypesCommand) listTypes() {
	fmt.Println("Listing types")
	knownTypes := findTypes()
	for _, t := range knownTypes {
		fmt.Println("    ", t)
	}
}

func findTypes() []string {
	types := make([]string, 0)

	attributes, err := os.Open(".gitattributes")
	if err != nil {
		return types // No .gitattibtues == no file types
	}

	scanner := bufio.NewScanner(attributes)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.Contains(line, "filter=media") {
			fields := strings.Fields(line)
			types = append(types, fields[0])
		}
	}

	return types
}

func init() {
	registerCommand("types", func(c *Command) RunnableCommand {
		return &TypesCommand{Command: c}
	})
}