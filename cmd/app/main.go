package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/omelete/sofredor-orchestrator/internal/app"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Criar instância da aplicação
	application := app.NewApp()

	// Configurar opções do Wails
	err := wails.Run(&options.App{
		Title:  "Sofredor Orchestrator",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        application.Startup,
		OnShutdown:       application.Shutdown,
		Bind: []interface{}{
			application,
		},
	})

	if err != nil {
		log.Fatal("Erro ao iniciar aplicação:", err)
	}
}
