// Package platform is reserved for shared kernel helpers extracted from
// side-car domains (watch, vdiupload), especially Win32 primitives.
//
// Planned layout (not moved yet — see docs/domains/README.md):
//
//	internal/platform/win32  — EnumWindows, focus, clipboard, SendInput
//
// Domains may depend on platform; platform must not depend on domains.
package platform
