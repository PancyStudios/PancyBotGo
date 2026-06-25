package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func main() {
	files := []string{
		"internal/commands/register.go",
		"internal/commands/utils/register.go",
		"internal/commands/music.go",
		"internal/commands/embeds/register.go",
		"internal/commands/mod/register.go",
		"internal/commands/premium/register.go",
		"internal/commands/config/register.go",
		"internal/commands/ia/register.go",
		"internal/commands/fun/register.go",
		"internal/commands/reactions/register.go",
		"internal/commands/economy/register.go",
		"internal/commands/shop/register.go",
		"internal/commands/security/register.go",
		"internal/commands/levels/register.go",
		"internal/commands/help/register.go",
	}

	for _, f := range files {
		content, _ := ioutil.ReadFile(f)
		re := regexp.MustCompile(`AddGlobalCommand|RegisterCommand`)
		matches := re.FindAllString(string(content), -1)
		if len(matches) > 0 {
			fmt.Printf("%s: %d calls\n", f, len(matches))
		}
	}
}
