package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var categoryEmojis = map[string]string{
	"config":   "⚙️",
	"dev":      "💻",
	"economy":  "💰",
	"embeds":   "📝",
	"fun":      "🎉",
	"ia":       "🤖",
	"levels":   "🌟",
	"mod":      "🛡️",
	"premium":  "💎",
	"reaction": "🎭",
	"security": "🔒",
	"utils":    "🧰",
}

func main() {
	root := "/home/turbis/GolandProjects/PancyBotGo/internal/commands"
	
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		parts := strings.Split(filepath.ToSlash(path), "/")
		if len(parts) < 2 {
			return nil
		}
		
		category := parts[len(parts)-2]
		emoji, ok := categoryEmojis[category]
		if !ok {
			return nil
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(contentBytes)

		// Regular expression to find Description: "..."
		// We only want to replace descriptions that do not already start with an emoji and a pipe
		// Usually they start with "
		
		re := regexp.MustCompile(`Description:\s*"([^"]+)"`)
		
		newContent := re.ReplaceAllStringFunc(content, func(match string) string {
			sub := re.FindStringSubmatch(match)
			if len(sub) > 1 {
				descText := sub[1]
				// If it already contains "|", it probably has an emoji
				if strings.Contains(descText, "|") {
					return match
				}
				// Otherwise add the emoji
				return fmt.Sprintf(`Description: "%s | %s"`, emoji, descText)
			}
			return match
		})

		if newContent != content {
			err = os.WriteFile(path, []byte(newContent), 0644)
			if err != nil {
				fmt.Printf("Error writing %s: %v\n", path, err)
			} else {
				fmt.Printf("Updated %s\n", path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directories: %v\n", err)
	}
}
