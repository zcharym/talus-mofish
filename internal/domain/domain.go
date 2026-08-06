// Package domain defines logical product domains (bounded contexts) for Talus Echo.
//
// Domains are feature areas with clear ownership boundaries. Desktop-agent domains
// plug into Management routes and Agent tools; side-car domains ship as separate
// CLIs / cloud workers and share only a thin platform kernel.
//
// Layout convention (incremental DDD):
//
//	internal/<domain>/     — domain code (watch, vdiupload, …)
//	docs/domains/<domain>/ — design docs and ownership notes
//	cmd/<cli>/             — delivery adapters for side-car domains
//	cloud/<service>/       — cloud adapters when needed
//
// Shared kernel packages stay flat under internal/ (appservice, agent, aiclient,
// auth, config, database, store, …) and must not import domain packages.
package domain

const (
	// English is the English Learning domain: vocabulary, reading, SRS, Anki import.
	// Route prefix: "english.*". Code currently lives across content/, store/, appservice/.
	English = "english"

	// Watch is Echo Watch: local VDI client screen OCR + rule alerts + Cloudflare push.
	// Delivery: cmd/echo-watch, cloud/echo-watch. Code: internal/watch.
	Watch = "watch"

	// VDIUpload is automated host→guest file transfer for H3C Workspace (and later RDP/Citrix/Horizon).
	// Delivery: cmd/vdi-upload. Code: internal/vdiupload. Design: docs/domains/vdiupload/DESIGN.md.
	VDIUpload = "vdiupload"
)

// All lists known domain identifiers for discovery and tooling.
var All = []string{English, Watch, VDIUpload}

// Info describes a domain for docs and tooling.
type Info struct {
	ID          string
	Name        string
	Kind        string // "desktop-agent" | "sidecar"
	CodeRoot    string
	DocsRoot    string
	Description string
}

// Catalog returns static metadata for each bounded context.
func Catalog() []Info {
	return []Info{
		{
			ID:          English,
			Name:        "English Learning",
			Kind:        "desktop-agent",
			CodeRoot:    "internal/content, internal/store, internal/appservice (english.*)",
			DocsRoot:    "docs/domains/english, docs/design-and-plan.md, docs/chat-learning-flows-plan.md",
			Description: "Vocabulary, articles, Anki import, SRS; first Wails Management/Agent domain.",
		},
		{
			ID:          Watch,
			Name:        "Echo Watch",
			Kind:        "sidecar",
			CodeRoot:    "internal/watch, cmd/echo-watch, cloud/echo-watch",
			DocsRoot:    "docs/domains/watch, cloud/echo-watch/README.md",
			Description: "Physical-PC screen watch of VDI client windows; OCR rules → Cloudflare → iOS PWA push.",
		},
		{
			ID:          VDIUpload,
			Name:        "VDI File Upload",
			Kind:        "sidecar",
			CodeRoot:    "internal/vdiupload, cmd/vdi-upload",
			DocsRoot:    "docs/domains/vdiupload",
			Description: "Automate file upload into H3C Workspace VDI via clipboard + Win32 UI automation.",
		},
	}
}
