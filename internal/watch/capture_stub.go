//go:build !windows

package watch

import "fmt"

type StubCapturer struct{}

func NewStubCapturer() *StubCapturer {
	return &StubCapturer{}
}

func (c *StubCapturer) Resolve(target Target) (WindowState, error) {
	return WindowState{Capturable: false, Reason: "platform_not_supported"}, nil
}

func (c *StubCapturer) CaptureRegion(state WindowState, region Region) ([]byte, error) {
	return nil, fmt.Errorf("screen capture is only supported on windows")
}

func (c *StubCapturer) OCR(pixels []byte, width, height int) (string, error) {
	return "", fmt.Errorf("ocr is only supported on windows")
}
