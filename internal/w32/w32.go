package w32

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                       = windows.NewLazySystemDLL("user32")
	procSendMessage              = user32.NewProc("SendMessageW")
	procMessageBox               = user32.NewProc("MessageBoxW")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procShowWindow               = user32.NewProc("ShowWindow")

	shell32                 = windows.NewLazySystemDLL("shell32")
	procExtractIcon         = shell32.NewProc("ExtractIconW")
	procSHBrowseForFolder   = shell32.NewProc("SHBrowseForFolderW")
	procSHGetPathFromIDList = shell32.NewProc("SHGetPathFromIDListW")

	kernel32             = windows.NewLazySystemDLL("kernel32")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")
)

// SHBrowseForFolder flags
const (
	BIF_RETURNONLYFSDIRS    = 0x00000001
	BIF_DONTGOBELOWDOMAIN   = 0x00000002
	BIF_STATUSTEXT          = 0x00000004
	BIF_RETURNFSANCESTORS   = 0x00000008
	BIF_EDITBOX             = 0x00000010
	BIF_VALIDATE            = 0x00000020
	BIF_NEWDIALOGSTYLE      = 0x00000040
	BIF_BROWSEINCLUDEURLS   = 0x00000080
	BIF_USENEWUI            = BIF_EDITBOX | BIF_NEWDIALOGSTYLE
	BIF_UAHINT              = 0x00000100
	BIF_NONEWFOLDERBUTTON   = 0x00000200
	BIF_NOTRANSLATETARGETS  = 0x00000400
	BIF_BROWSEFORCOMPUTER   = 0x00001000
	BIF_BROWSEFORPRINTER    = 0x00002000
	BIF_BROWSEINCLUDEFILES  = 0x00004000
	BIF_SHAREABLE           = 0x00008000
	BIF_BROWSEFILEJUNCTIONS = 0x00010000
)

type (
	DWORD  uint32
	HANDLE uintptr
	HWND   HANDLE
)

// BROWSEINFO http://msdn.microsoft.com/en-us/library/windows/desktop/bb773205.aspx
type BROWSEINFO struct {
	Owner        HWND
	Root         *uint16
	DisplayName  *uint16
	Title        *uint16
	Flags        uint32
	CallbackFunc uintptr
	LParam       uintptr
	Image        int32
}

func ExtractIcon(exeFileName string, iconIndex int32) uintptr {
	e, _ := syscall.UTF16PtrFromString(exeFileName)
	ret, _, _ := procExtractIcon.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(e)),
		uintptr(iconIndex),
	)
	return ret
}

func SendMessage(hwnd unsafe.Pointer, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := syscall.Syscall6(procSendMessage.Addr(), 4,
		uintptr(hwnd), uintptr(msg),
		wParam, lParam,
		0, 0)
	return ret
}

func MessageBox(hwnd unsafe.Pointer, text, caption string, flags uint) uintptr {
	_text, _ := syscall.UTF16PtrFromString(text)
	_caption, _ := syscall.UTF16PtrFromString(caption)
	ret, _, _ := syscall.Syscall6(procMessageBox.Addr(), 4,
		uintptr(hwnd),
		uintptr(unsafe.Pointer(_text)),
		uintptr(unsafe.Pointer(_caption)),
		uintptr(flags),
		0, 0,
	)
	return ret
}

func ProcessExistsByProcessName(name string) (exists bool, err error) {
	cmd := exec.Command("cmd", "/C", "tasklist", "|", "findstr", name)
	out, err := cmd.Output()
	if err != nil {
		return
	}
	exists = len(strings.Fields(string(out))) > 0
	return
}

func GetConsoleWindow() uintptr {
	ret, _, _ := procGetConsoleWindow.Call()
	return ret
}

func GetWindowThreadProcessId(hwnd uintptr) (uintptr, int) {
	var pid int
	ret, _, _ := procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return ret, pid
}

func ShowWindow(hwnd uintptr, show int) bool {
	ret, _, _ := procShowWindow.Call(hwnd, uintptr(show))
	return ret != 0
}

func HideConsoleWindow() {
	hwnd := GetConsoleWindow()
	if hwnd <= 0 {
		return
	}
	_, pid := GetWindowThreadProcessId(hwnd)
	if pid == os.Getpid() {
		ShowWindow(hwnd, 0)
	}
}

func SHGetPathFromIDList(idl uintptr) string {
	buf := make([]uint16, 1024)
	_, _, _ = procSHGetPathFromIDList.Call(
		idl,
		uintptr(unsafe.Pointer(&buf[0])),
	)
	return syscall.UTF16ToString(buf)
}

func SHBrowseForFolder(bi *BROWSEINFO) uintptr {
	ret, _, _ := procSHBrowseForFolder.Call(uintptr(unsafe.Pointer(bi)))
	return ret
}
