package web

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

// localOnlyMiddleware ensures a route is only accessible via localhost
func localOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		if !strings.HasPrefix(host, "localhost") && !strings.HasPrefix(host, "127.0.0.1") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "El panel de administración solo es accesible localmente."})
			return
		}
		c.Next()
	}
}

// SetupDashboardRoutes configures the local admin dashboard endpoints
func SetupDashboardRoutes(s *Server) {
	// Group admin routes with the local-only middleware
	admin := s.Group("/admin", localOnlyMiddleware())

	// Serve the main dashboard HTML
	admin.GET("", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(dashboardHTML))
	})

	// API endpoints for the dashboard
	api := admin.Group("/api")

	api.GET("/stats", func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		c.JSON(http.StatusOK, gin.H{
			"goroutines": runtime.NumGoroutine(),
			"memory_mb":  float64(m.Alloc) / 1024 / 1024,
			"sys_mb":     float64(m.Sys) / 1024 / 1024,
			"status":     "online",
		})
	})

	api.POST("/stop", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Apagando..."})

		// Send signal to main to gracefully shutdown
		go func() {
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()
	})

	api.POST("/command", func(c *gin.Context) {
		var req struct {
			Command string `json:"command"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		cmd := strings.TrimSpace(strings.ToLower(req.Command))
		var output string

		switch cmd {
		case "help":
			output = "=== PancyBot Web CLI ===\n  help  - Muestra esta lista\n  stats - Estadísticas de sistema\n  clear - Limpia la consola\n  stop  - Apaga el bot\n========================"
		case "stats":
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			output = fmt.Sprintf("=== Estadísticas ===\nGoroutines: %d\nMemoria OS: %.2f MB\nMemoria Uso: %.2f MB\nGC Pausas: %d\n===================",
				runtime.NumGoroutine(), float64(m.Sys)/1024/1024, float64(m.Alloc)/1024/1024, m.NumGC)
		case "clear":
			output = "_CLEAR_"
		case "stop", "exit", "quit":
			output = "Iniciando apagado seguro..."
			go func() {
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}()
		case "":
			output = ""
		default:
			output = fmt.Sprintf("Comando desconocido: '%s'. Escribe 'help'.", cmd)
		}

		c.JSON(http.StatusOK, gin.H{"output": output})
	})
}

const dashboardHTML = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>PancyBot | Admin Dashboard</title>
    <style>
        :root {
            --bg: #0f172a;
            --surface: #1e293b;
            --primary: #3b82f6;
            --primary-hover: #2563eb;
            --text: #f8fafc;
            --text-muted: #94a3b8;
            --danger: #ef4444;
        }
        body {
            font-family: 'Inter', -apple-system, sans-serif;
            background-color: var(--bg);
            color: var(--text);
            margin: 0;
            padding: 2rem;
            display: flex;
            flex-direction: column;
            align-items: center;
        }
        .container {
            width: 100%;
            max-width: 800px;
        }
        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 2rem;
            padding-bottom: 1rem;
            border-bottom: 1px solid var(--surface);
        }
        .title {
            font-size: 1.5rem;
            font-weight: bold;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .dot {
            width: 12px;
            height: 12px;
            background-color: #22c55e;
            border-radius: 50%;
            box-shadow: 0 0 10px #22c55e;
        }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1rem;
            margin-bottom: 2rem;
        }
        .card {
            background-color: var(--surface);
            padding: 1.5rem;
            border-radius: 12px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
        }
        .card-title {
            font-size: 0.875rem;
            color: var(--text-muted);
            margin-bottom: 0.5rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }
        .card-value {
            font-size: 2rem;
            font-weight: bold;
        }
        .actions {
            display: flex;
            gap: 1rem;
        }
        .btn {
            background-color: var(--primary);
            color: white;
            border: none;
            padding: 0.75rem 1.5rem;
            border-radius: 8px;
            font-size: 1rem;
            font-weight: 500;
            cursor: pointer;
            transition: background-color 0.2s;
        }
        .btn:hover { background-color: var(--primary-hover); }
        .btn-danger { background-color: var(--danger); }
        .btn-danger:hover { background-color: #dc2626; }
        
        /* Terminal Styles */
        .terminal-container {
            width: 100%;
            background-color: #000;
            border-radius: 8px;
            overflow: hidden;
            margin-top: 2rem;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.5);
            display: flex;
            flex-direction: column;
            border: 1px solid var(--surface);
        }
        .terminal-header {
            background-color: #1a1a1a;
            padding: 8px 16px;
            font-size: 0.8rem;
            color: #888;
            font-family: monospace;
            border-bottom: 1px solid #333;
            display: flex;
            gap: 6px;
        }
        .terminal-dot { width: 10px; height: 10px; border-radius: 50%; }
        .terminal-dot.r { background-color: #ff5f56; }
        .terminal-dot.y { background-color: #ffbd2e; }
        .terminal-dot.g { background-color: #27c93f; }
        
        .terminal-output {
            padding: 16px;
            height: 250px;
            overflow-y: auto;
            font-family: 'Consolas', 'Monaco', monospace;
            font-size: 0.9rem;
            color: #0f0;
            white-space: pre-wrap;
            line-height: 1.4;
        }
        .terminal-input-row {
            display: flex;
            padding: 0 16px 16px 16px;
            font-family: 'Consolas', 'Monaco', monospace;
            color: #0f0;
            align-items: center;
        }
        .terminal-prompt { margin-right: 8px; font-weight: bold; }
        .terminal-input {
            flex: 1;
            background: transparent;
            border: none;
            color: #fff;
            font-family: 'Consolas', 'Monaco', monospace;
            font-size: 0.9rem;
            outline: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="title">
                <div class="dot"></div>
                PancyBot Admin
            </div>
            <div>Local Dashboard</div>
        </div>

        <div class="grid">
            <div class="card">
                <div class="card-title">Memoria RAM (Uso)</div>
                <div class="card-value" id="mem-alloc">-- MB</div>
            </div>
            <div class="card">
                <div class="card-title">Memoria Sistema</div>
                <div class="card-value" id="mem-sys">-- MB</div>
            </div>
            <div class="card">
                <div class="card-title">Goroutines</div>
                <div class="card-value" id="goroutines">--</div>
            </div>
        </div>

        <div class="actions">
            <button class="btn" onclick="fetchStats()">Refrescar Datos</button>
            <button class="btn btn-danger" onclick="stopBot()">Apagar Bot</button>
        </div>

        <div class="terminal-container">
            <div class="terminal-header">
                <div class="terminal-dot r"></div>
                <div class="terminal-dot y"></div>
                <div class="terminal-dot g"></div>
                <span style="margin-left:10px">PancyBot Web Console</span>
            </div>
            <div class="terminal-output" id="term-out">PancyBot Go [Versión 1.0.0]
(c) Pancy Studios. Todos los derechos reservados.

Escribe 'help' para ver los comandos disponibles.
</div>
            <div class="terminal-input-row">
                <span class="terminal-prompt">admin@pancybot:~$</span>
                <input type="text" class="terminal-input" id="term-in" autocomplete="off" spellcheck="false" autofocus>
            </div>
        </div>
    </div>

    <script>
        async function fetchStats() {
            try {
                const res = await fetch('/admin/api/stats');
                const data = await res.json();
                document.getElementById('mem-alloc').innerText = data.memory_mb.toFixed(2) + ' MB';
                document.getElementById('mem-sys').innerText = data.sys_mb.toFixed(2) + ' MB';
                document.getElementById('goroutines').innerText = data.goroutines;
            } catch (e) {
                console.error(e);
            }
        }

        async function stopBot() {
            if(confirm("¿Estás seguro de que quieres apagar el bot?")) {
                try {
                    await fetch('/admin/api/stop', {method: 'POST'});
                    alert("Se ha enviado la señal de apagado. Puedes cerrar esta ventana.");
                } catch(e) {
                    alert("Error al apagar");
                }
            }
        }

        // Fetch on load
        fetchStats();
        // Auto refresh every 5s
        setInterval(fetchStats, 5000);

        // Terminal logic
        const termIn = document.getElementById('term-in');
        const termOut = document.getElementById('term-out');

        function appendTerm(text) {
            termOut.textContent += '\n' + text;
            termOut.scrollTop = termOut.scrollHeight;
        }

        termIn.addEventListener('keypress', async function(e) {
            if (e.key === 'Enter') {
                const cmd = termIn.value.trim();
                termIn.value = '';
                
                appendTerm('admin@pancybot:~$ ' + cmd);
                
                if (cmd === '') return;

                try {
                    const res = await fetch('/admin/api/command', {
                        method: 'POST',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify({command: cmd})
                    });
                    const data = await res.json();
                    
                    if (data.output === '_CLEAR_') {
                        termOut.textContent = '';
                    } else if (data.output) {
                        appendTerm(data.output);
                    }
                } catch(err) {
                    appendTerm('Error de conexión con el servidor local.');
                }
            }
        });
        
        // Focus terminal when clicking on it
        document.querySelector('.terminal-container').addEventListener('click', () => {
            termIn.focus();
        });
    </script>
</body>
</html>
`
