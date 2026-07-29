//go:build windows

package watch

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	gdi32                    = windows.NewLazySystemDLL("gdi32.dll")
	procGetWindowDC          = user32.NewProc("GetWindowDC")
	procReleaseDC            = user32.NewProc("ReleaseDC")
	procBitBlt               = gdi32.NewProc("BitBlt")
	procCreateCompatibleDC   = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procSelectObject         = gdi32.NewProc("SelectObject")
	procDeleteObject         = gdi32.NewProc("DeleteObject")
	procDeleteDC             = gdi32.NewProc("DeleteDC")
	procGetDIBits            = gdi32.NewProc("GetDIBits")
)

const srccopy = 0x00CC0020

type WindowsCapturer struct{}

func NewWindowsCapturer() *WindowsCapturer {
	return &WindowsCapturer{}
}

func (c *WindowsCapturer) Resolve(target Target) (WindowState, error) {
	return resolveWindowState(target)
}

func (c *WindowsCapturer) CaptureRegion(state WindowState, region Region) ([]byte, error) {
	if !state.Capturable || state.HWND == 0 {
		return nil, fmt.Errorf("window not capturable: %s", state.Reason)
	}
	if region.W <= 0 || region.H <= 0 {
		return nil, fmt.Errorf("invalid region")
	}

	hdcWindow, _, err := procGetWindowDC.Call(state.HWND)
	if hdcWindow == 0 {
		return nil, fmt.Errorf("GetWindowDC: %w", err)
	}
	defer procReleaseDC.Call(state.HWND, hdcWindow)

	hdcMem, _, err := procCreateCompatibleDC.Call(hdcWindow)
	if hdcMem == 0 {
		return nil, fmt.Errorf("CreateCompatibleDC: %w", err)
	}
	defer procDeleteDC.Call(hdcMem)

	hBitmap, _, err := procCreateCompatibleBitmap.Call(hdcWindow, uintptr(region.W), uintptr(region.H))
	if hBitmap == 0 {
		return nil, fmt.Errorf("CreateCompatibleBitmap: %w", err)
	}
	defer procDeleteObject.Call(hBitmap)

	procSelectObject.Call(hdcMem, hBitmap)
	ok, _, err := procBitBlt.Call(
		hdcMem, 0, 0, uintptr(region.W), uintptr(region.H),
		hdcWindow, uintptr(region.X), uintptr(region.Y), srccopy,
	)
	if ok == 0 {
		return nil, fmt.Errorf("BitBlt: %w", err)
	}

	return bitmapToRGBA(hdcMem, hBitmap, region.W, region.H)
}

func (c *WindowsCapturer) OCR(pixels []byte, width, height int) (string, error) {
	img := rgbaToImage(pixels, width, height)
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	tmp, err := os.CreateTemp("", "echo-watch-*.png")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	defer os.Remove(path)

	if _, err := tmp.Write(buf.Bytes()); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}

	script := fmt.Sprintf(ocrScript, filepath.ToSlash(path))
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ocr powershell: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

const ocrScript = `
Add-Type -AssemblyName System.Runtime.WindowsRuntime
$null = [Windows.Storage.StorageFile, Windows.Storage, ContentType=WindowsRuntime]
$null = [Windows.Graphics.Imaging.BitmapDecoder, Windows.Graphics.Imaging, ContentType=WindowsRuntime]
$null = [Windows.Media.Ocr.OcrEngine, Windows.Media.Ocr, ContentType=WindowsRuntime]
$path = "%s"
$asyncOp = [Windows.Storage.StorageFile]::GetFileFromPathAsync($path)
$file = $asyncOp.GetAwaiter().GetResult()
$streamOp = $file.OpenAsync([Windows.Storage.FileAccessMode]::Read)
$stream = $streamOp.GetAwaiter().GetResult()
$decoderOp = [Windows.Graphics.Imaging.BitmapDecoder]::CreateAsync($stream)
$decoder = $decoderOp.GetAwaiter().GetResult()
$bitmapOp = $decoder.GetSoftwareBitmapAsync()
$bitmap = $bitmapOp.GetAwaiter().GetResult()
$engine = [Windows.Media.Ocr.OcrEngine]::TryCreateFromUserProfileLanguages()
$resultOp = $engine.RecognizeAsync($bitmap)
$result = $resultOp.GetAwaiter().GetResult()
Write-Output $result.Text
`

func bitmapToRGBA(hdcMem, hBitmap uintptr, width, height int) ([]byte, error) {
	var bi struct {
		Size          uint32
		Width         int32
		Height        int32
		Planes        uint16
		BitCount      uint16
		Compression   uint32
		SizeImage     uint32
		XPelsPerMeter int32
		YPelsPerMeter int32
		ClrUsed       uint32
		ClrImportant  uint32
	}
	bi.Size = uint32(unsafe.Sizeof(bi))
	bi.Width = int32(width)
	bi.Height = -int32(height)
	bi.Planes = 1
	bi.BitCount = 32

	buf := make([]byte, width*height*4)
	ret, _, err := procGetDIBits.Call(
		hdcMem, hBitmap, 0, uintptr(height),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bi)), 0,
	)
	if ret == 0 {
		return nil, fmt.Errorf("GetDIBits: %w", err)
	}
	return buf, nil
}

func rgbaToImage(pixels []byte, width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 4
			img.SetRGBA(x, y, color.RGBA{
				R: pixels[i+2],
				G: pixels[i+1],
				B: pixels[i+0],
				A: 255,
			})
		}
	}
	return img
}
