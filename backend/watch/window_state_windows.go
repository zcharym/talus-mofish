//go:build windows

package watch

import (
	"fmt"
	"regexp"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	kernel32             = windows.NewLazySystemDLL("kernel32.dll")
	procEnumWindows      = user32.NewProc("EnumWindows")
	procGetWindowTextW   = user32.NewProc("GetWindowTextW")
	procGetWindowThread  = user32.NewProc("GetWindowThreadProcessId")
	procIsIconic         = user32.NewProc("IsIconic")
	procIsWindowVisible  = user32.NewProc("IsWindowVisible")
	procGetWindowRect    = user32.NewProc("GetWindowRect")
	procCreateSnapshot   = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW  = kernel32.NewProc("Process32FirstW")
	procProcess32NextW   = kernel32.NewProc("Process32NextW")
	procCloseHandle      = kernel32.NewProc("CloseHandle")
)

const (
	th32csSnapProcess = 0x00000002
	maxPath           = 260
)

type processEntry32 struct {
	Size            uint32
	Usage           uint32
	ProcessID       uint32
	DefaultHeapID   uintptr
	ModuleID        uint32
	Threads         uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	ExeFile         [maxPath]uint16
}

type rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type windowSearch struct {
	titlePattern *regexp.Regexp
	processName  string
	processIDs   map[uint32]struct{}
	bestHWND     uintptr
}

func processRunning(name string) (map[uint32]struct{}, error) {
	if name == "" {
		return nil, fmt.Errorf("process_name is required")
	}
	name = strings.ToLower(name)

	snap, _, err := procCreateSnapshot.Call(th32csSnapProcess, 0)
	if snap == uintptr(windows.InvalidHandle) {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot: %w", err)
	}
	defer procCloseHandle.Call(snap)

	var entry processEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	ret, _, err := procProcess32FirstW.Call(snap, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return nil, fmt.Errorf("Process32First: %w", err)
	}

	ids := make(map[uint32]struct{})
	for {
		exe := strings.ToLower(windows.UTF16ToString(entry.ExeFile[:]))
		if exe == name {
			ids[entry.ProcessID] = struct{}{}
		}
		ret, _, err = procProcess32NextW.Call(snap, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("vdi_process_not_found")
	}
	return ids, nil
}

func enumWindowCallback(hwnd uintptr, lparam uintptr) uintptr {
	search := (*windowSearch)(unsafe.Pointer(lparam))

	visible, _, _ := procIsWindowVisible.Call(hwnd)
	if visible == 0 {
		return 1
	}

	var pid uint32
	procGetWindowThread.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if len(search.processIDs) > 0 {
		if _, ok := search.processIDs[pid]; !ok {
			return 1
		}
	}

	title := windowTitle(hwnd)
	if title == "" {
		return 1
	}
	if search.titlePattern != nil && !search.titlePattern.MatchString(title) {
		return 1
	}

	search.bestHWND = hwnd
	return 0
}

func windowTitle(hwnd uintptr) string {
	buf := make([]uint16, 512)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return windows.UTF16ToString(buf)
}

func findClientWindow(target Target) (uintptr, error) {
	pattern, err := target.TitlePattern()
	if err != nil {
		return 0, err
	}

	ids, err := processRunning(target.ProcessName)
	if err != nil {
		return 0, err
	}

	search := &windowSearch{
		titlePattern: pattern,
		processName:  target.ProcessName,
		processIDs:   ids,
	}
	procEnumWindows.Call(
		syscall.NewCallback(enumWindowCallback),
		uintptr(unsafe.Pointer(search)),
	)
	if search.bestHWND == 0 {
		return 0, fmt.Errorf("vdi_window_not_found")
	}
	return search.bestHWND, nil
}

func inspectWindowState(hwnd uintptr, requireVisible bool) WindowState {
	if hwnd == 0 {
		return WindowState{Capturable: false, Reason: "vdi_window_not_found"}
	}

	iconic, _, _ := procIsIconic.Call(hwnd)
	if iconic != 0 {
		return WindowState{Capturable: false, Reason: "vdi_window_minimized", HWND: hwnd}
	}

	visible, _, _ := procIsWindowVisible.Call(hwnd)
	if requireVisible && visible == 0 {
		return WindowState{Capturable: false, Reason: "vdi_window_not_visible", HWND: hwnd}
	}

	var r rect
	ok, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	if ok == 0 {
		return WindowState{Capturable: false, Reason: "vdi_window_not_visible", HWND: hwnd}
	}
	width := int(r.Right - r.Left)
	height := int(r.Bottom - r.Top)
	if width <= 0 || height <= 0 {
		return WindowState{Capturable: false, Reason: "vdi_window_not_visible", HWND: hwnd}
	}

	return WindowState{
		Capturable: true,
		Reason:     "ok",
		HWND:       hwnd,
		Width:      width,
		Height:     height,
	}
}

func resolveWindowState(target Target) (WindowState, error) {
	hwnd, err := findClientWindow(target)
	if err != nil {
		return WindowState{Capturable: false, Reason: err.Error()}, nil
	}
	return inspectWindowState(hwnd, target.RequireVisibleWindow()), nil
}
