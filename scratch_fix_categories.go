package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	baseDir := "./internal/messagecommands"
	
	filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "register.go" {
			// Extract category from folder name
			dir := filepath.Dir(path)
			category := filepath.Base(dir)
			
			// Capitalize category name
			if len(category) > 0 {
				category = strings.ToUpper(category[:1]) + category[1:]
			}

			contentBytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			content := string(contentBytes)
			
			// We only want to replace `"General"` with `"` + category + `"` on the lines with RegisterCommand
			lines := strings.Split(content, "\n")
			changed := false
			for i, line := range lines {
				if strings.Contains(line, "messagecommands.RegisterCommand") {
					// We find "General" and replace it
					if strings.Contains(line, `"General"`) {
						lines[i] = strings.Replace(line, `"General"`, `"`+category+`"`, 1)
						changed = true
					}
				}
			}
			
			if changed {
				newContent := strings.Join(lines, "\n")
				err = os.WriteFile(path, []byte(newContent), 0644)
				if err != nil {
					fmt.Println("Error writing", path, ":", err)
				} else {
					fmt.Println("Updated category in", path, "to", category)
				}
			}
		}
		return nil
	})
}
