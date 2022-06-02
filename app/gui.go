package app

import (
	"github.com/jchv/go-webview2"
	"github.com/jchv/go-webview2/pkg/edge"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"
	"vrchat-osc-manager/internal/w32"
)

const guiTitle = "VRChat OSC Manager"

var (
	ui webview2.WebView
)

func GUI() {
	if RuntimeVersion() == "" {
		w32.MessageBox(0, "Please install Webview2 Runtime to use this application.", guiTitle, 0)
		return
	}
	ui = webview2.NewWithOptions(webview2.WebViewOptions{
		WindowOptions: webview2.WindowOptions{Title: guiTitle},
	})
	defer ui.Destroy()

	w32.SendMessage(ui.Window(), 0x0080, 1, w32.ExtractIcon(os.Args[0], 0))

	ui.SetSize(800, 600, webview2.HintNone)

	chromium := GetChromium(ui)
	settings, _ := chromium.GetSettings()
	if !*debugMode {
		_ = settings.PutAreDevToolsEnabled(false)
	}
	_ = settings.PutIsStatusBarEnabled(false)

	folderPath, _ := filepath.Abs("./public")
	webview := chromium.GetICoreWebView2_3()
	_ = webview.SetVirtualHostNameToFolderMapping(
		"app.assets", folderPath,
		edge.COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND_DENY_CORS,
	)

	ui.Navigate("http://app.assets/index.html")
	ui.Run()
}

func GetChromium(w webview2.WebView) *edge.Chromium {
	browser := reflect.ValueOf(w).Elem().FieldByName("browser")
	browser = reflect.NewAt(browser.Type(), unsafe.Pointer(browser.UnsafeAddr())).Elem()
	chromium, ok := browser.Interface().(*edge.Chromium)
	if ok {
		return chromium
	}
	return nil
}

func RuntimeVersion() (version string) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\WOW6432Node\\Microsoft\\EdgeUpdate\\Clients\\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}", registry.READ)
	if err != nil {
		k, err = registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\EdgeUpdate\\Clients\\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}", registry.READ)
		if err != nil {
			return
		}
	}
	version, _, _ = k.GetStringValue("pv")
	return
}
