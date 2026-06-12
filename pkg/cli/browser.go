package cli

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/PancyStudios/PancyBotGo/pkg/logger"
)

// OpenBrowser attempts to open the given URL in the user's default browser
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		logger.Warn(fmt.Sprintf("No se pudo abrir el navegador automáticamente: %v. Por favor, visita %s manualmente.", err, url), "CLI")
	} else {
		logger.Info(fmt.Sprintf("Navegador abierto en %s", url), "CLI")
	}
}
