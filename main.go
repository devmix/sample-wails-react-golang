// main initializes and runs the ToDo Notes desktop application built with Wails.
// It sets up the SQLite database, creates repository and service layers,
// configures the application window, and starts the event loop.
package main

import (
	"embed"
	"log"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/wailsapp/wails/v3/pkg/application"

	"sample-wails-react-golang/internal/app/service"
	"sample-wails-react-golang/internal/pkg/repository"
)

//go:embed all:frontend/dist
var assets embed.FS

// main is the entry point for the application. It configures the database path
// under the XDG data directory, initializes repositories and services, creates
// the Wails application with its window settings, and starts the runtime.
func main() {
	dbPath := filepath.Join(xdg.DataHome, "sample-wails-react-golang", "app.db")

	db, err := repository.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	todoRepo := repository.NewTodoRepo(db)
	noteRepo := repository.NewNoteRepo(db)
	tagRepo := repository.NewTagRepo(db)

	service.SetVersion("0.0.1")

	app := application.New(application.Options{
		Name:        "ToDo Notes",
		Description: "A to-do and notes application",
		Services: []application.Service{
			application.NewService(service.NewTodoService(todoRepo)),
			application.NewService(service.NewNoteService(noteRepo)),
			application.NewService(service.NewTagService(tagRepo)),
			application.NewService(service.NewAppService()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "ToDo Notes",
		Width:  1200,
		Height: 800,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		URL: "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
