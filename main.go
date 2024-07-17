package main

import (
	"context"
	"embed"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"net/http/pprof"
	"terminal/logic"
	"time"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	logicApp := logic.NewApp()
	app := application.NewWithOptions(&options.App{
		Title:     "terminal",
		Width:     924,
		Height:    520,
		MinWidth:  800,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        true,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup: func(ctx context.Context) {
			runtime.WindowSetDarkTheme(ctx)
			logicApp.Ctx = ctx
		},
		Bind: []interface{}{
			logicApp,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			BackdropType:                      windows.Mica,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.Dark,
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
			WindowIsTranslucent: true,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways,
			ProgramName:         "terminal-console",
		},
	})
	//go startPprof()
	err := app.Run()
	if err != nil {
		println("Error:", err.Error())
	}
}

func startPprof() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.HideBanner = true
	RegisterRoutes(e)
	e.Logger.Fatal(e.StartServer(&http.Server{
		Addr:              ":1323",
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 2,
		WriteTimeout:      time.Second * 120,
	}))
}
func RegisterRoutes(engine *echo.Echo) {
	router := engine.Group("/debug")
	// 下面的路由根据要采集的数据需求注册，不用全都注册
	router.GET("/pprof", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/pprof/allocs", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/pprof/block", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/pprof/goroutine", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/pprof/heap", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/pprof/mutex", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
	router.GET("/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	router.GET("/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	router.GET("/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	router.GET("/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
}
