//go:build tools

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

func main() {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := stealth.MustPage(browser)

	fmt.Println("Navigating to craiyon.com...")
	page.MustNavigate("https://www.craiyon.com")

	// Wait until the Cloudflare challenge is passed
	page.MustWaitLoad()
	fmt.Println("Page loaded. Waiting 5s...")
	time.Sleep(5 * time.Second)

	fmt.Println("Running fetch for API...")

	js := fmt.Sprintf(`
		async () => {
			let res = await fetch("https://api.craiyon.com/v4", {
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify({
					prompt: "A beautiful red dog",
					token: %q,
					model: "auto",
					negative_prompt: "",
					size: "256x256"
				})
			});
			return await res.json();
		}
	`, os.Getenv("CRAIYON_TEST_TOKEN"))

	res, err := page.Eval(js)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Result: %v\n", res.Value.String())
}
