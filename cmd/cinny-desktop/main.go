package main

import (
	"fmt"
	"github.com/diamondburned/gotk4-webkitgtk/pkg/webkit2/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
	"github.com/yinzhu32/cinny-desktop/pkg/assets/cinny"
	"github.com/yinzhu32/cinny-desktop/pkg/assets/glade"
	"net/http"
	"os"
	"time"
)

const (
	serverProtocol = "https"
	serverHost     = "app.element.io"
	serverPort     = 443
	serverLocal    = false
)

func main() {
	dbusConn, err := dbus.SessionBusPrivate()
	if err != nil {
		panic(err)
	}

	app := gtk.NewApplication("com.github.yinzhu32.cinny-desktop", 0)
	app.Connect("activate", func() {
		builder := gtk.NewBuilderFromString(glade.MainApplicationWindow, len(glade.MainApplicationWindow))

		webView := builder.GetObject("webview").Cast().(*webkit2.WebView)
		webView.LoadURI(fmt.Sprintf("%s://%s:%d", serverProtocol, serverHost, serverPort))
		webView.Connect("show-notification", func(view *webkit2.WebView, notification *webkit2.Notification) {
			println("###### SHOW NOTIFICATION!!!! #####")
			println(notification.Title(), notification.Body())
			_, err := notify.SendNotification(dbusConn, notify.Notification{
				AppName:       "cinny-desktop",
				ReplacesID:    0,
				AppIcon:       "mail-unread",
				Summary:       notification.Title(),
				Body:          notification.Body(),
				Actions:       []notify.Action{
					{Key: "cancel", Label: "Cancel"},
				},
				Hints:         nil,
				ExpireTimeout: 5 * time.Second,
			})
			if err != nil {
				println("## NOTIFICATION NOT SENT ##", err)
				return
			}
		})
		webView.Connect("permission-request", func(request *webkit2.PermissionRequest) {
			println("############# PERMISSION REQUEST!!!!! ###############", request)
			request.Allow()
		})
		securityOrigin := webkit2.NewSecurityOrigin(serverProtocol, serverHost, uint16(serverPort))
		webView.Context().InitializeNotificationPermissions([]*webkit2.SecurityOrigin{securityOrigin}, nil)

		mainWindow := builder.GetObject("mainwindow").Cast().(*gtk.ApplicationWindow)
		mainWindow.SetApplication(app)
		mainWindow.SetTitle("Cinny Desktop")
		mainWindow.SetDefaultSize(800, 600)
		mainWindow.ShowAll()
	})

	go func() {
		if !serverLocal {
			return
		}
		mux := http.NewServeMux()
		mux.Handle("/", http.FileServer(http.FS(cinny.Filesystem)))
		appServer := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", serverHost, serverPort),
			Handler: mux,
		}
		if err := appServer.ListenAndServe(); err != nil {
			app.Quit()
			return
		}
	}()

	code := app.Run(os.Args)
	os.Exit(code)
}
