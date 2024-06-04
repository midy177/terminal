package main

import (
	"embed"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "terminal",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        true,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			BackdropType:                      windows.Mica,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				// Theme to use when window is active
				DarkModeTitleBar:   windows.RGB(255, 0, 0), // Red
				DarkModeTitleText:  windows.RGB(0, 255, 0), // Green
				DarkModeBorder:     windows.RGB(0, 0, 255), // Blue
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
				// Theme to use when window is inactive
				DarkModeTitleBarInactive:   windows.RGB(128, 0, 0),
				DarkModeTitleTextInactive:  windows.RGB(0, 128, 0),
				DarkModeBorderInactive:     windows.RGB(0, 0, 128),
				LightModeTitleBarInactive:  windows.RGB(100, 100, 100),
				LightModeTitleTextInactive: windows.RGB(10, 10, 10),
				LightModeBorderInactive:    windows.RGB(100, 100, 100),
			},
			// User messages that can be customised
			//Messages *windows.Messages
			// OnSuspend is called when Windows enters low power mode
			//OnSuspend func()
			// OnResume is called when Windows resumes from low power mode
			//OnResume func(),
			//WebviewGpuDisabled: false,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "Terminal Console",
				Message: "© 2024 Wuly",
				//Icon:    icon,
			},
		},
		Linux: &linux.Options{
			//Icon: icon,
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways,
			ProgramName:         "terminal-console",
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
