package main

import (
	"context"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os/exec"
	"time"
)

func setup() {
	var runtimeContext context.Context
	app := NewApp()

	appOptions := &options.App{
		Title:  app.Title(),
		Width:  700,
		Height: 520,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			runtimeContext = ctx
			runtime.EventsOn(runtimeContext, EventToggleOnLineStatus, func(optionalData ...interface{}) {
				runtime.WindowSetTitle(runtimeContext, app.Title())
			})

		},
		OnShutdown: func(ctx context.Context) {
			app.ShutdownN2N()
			<-time.After(time.Millisecond * 10)
		},
		Bind: []interface{}{
			app,
		},
	}

	err := wails.Run(appOptions)

	if err != nil {
		println("Error:", err.Error())
	}
}

func hideCmdWindow(c *exec.Cmd) {
	// do nothing
}
