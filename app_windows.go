package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os/exec"
	"syscall"
	"time"
)

func setup() {
	var runtimeContext context.Context
	var systray *application.SystemTray

	app := NewApp()
	var toggleConnectMenuItem *menu.MenuItem

	appOptions := &options.App{
		Title:             app.Title(),
		Width:             800,
		Height:            580,
		HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			runtimeContext = ctx

			runtime.EventsOn(runtimeContext, EventToggleOnLineStatus, func(optionalData ...interface{}) {
				if toggleConnectMenuItem != nil {
					isOnLine := app.IsOnline()
					toggleConnectMenuItem.SetChecked(isOnLine)
					runtime.WindowSetTitle(runtimeContext, app.Title())
					systray.Update()
				}
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

	mainApp := application.NewWithOptions(
		appOptions,
	)

	var showWindow = func() {
		runtime.WindowShow(runtimeContext)
	}

	systray = mainApp.NewSystemTray(&options.SystemTray{
		// This is the icon used when the system in using light mode
		LightModeIcon: &options.SystemTrayIcon{
			Data: lightModeIcon,
		},
		// This is the icon used when the system in using dark mode
		DarkModeIcon: &options.SystemTrayIcon{
			Data: darkModeIcon,
		},
		Tooltip:     "N2N User Edge",
		OnLeftClick: showWindow,
		OnMenuClose: func() {
			// Add the left click call after 500ms
			// We do this because the left click fires right
			// after the menu closes, and we don't want to show
			// the window on menu close.
			go func() {
				time.Sleep(500 * time.Millisecond)
				systray.OnLeftClick(showWindow)
			}()
		},
		OnMenuOpen: func() {
			// Remove the left click callback
			systray.OnLeftClick(func() {})
		},
	})

	showWindowMenuItem := menu.Label("Show window").OnClick(func(data *menu.CallbackData) {
		showWindow()
	})
	toggleConnectMenuItem = menu.Label("Online").SetChecked(app.IsOnline()).OnClick(func(data *menu.CallbackData) {
		if app.IsOnline() {
			app.ShutdownN2N()
		} else {
			app.SetupN2N()
		}
	})

	// Now we set the menu of the systray.
	// This would likely be created in a different function/file
	mainMenus := menu.NewMenuFromItems(
		showWindowMenuItem,
		menu.Separator(),
		toggleConnectMenuItem,
		menu.Separator(),
		menu.Label("Quit").OnClick(func(_ *menu.CallbackData) {
			println("Quitting application")
			mainApp.Quit()
		}),
	)
	systray.SetMenu(mainMenus)

	err := mainApp.Run()
	if err != nil {
		println("Error:", err.Error())
	}
}

func hideCmdWindow(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
