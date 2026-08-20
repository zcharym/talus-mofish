// Package consts defines logical product domains and shared constants for Talus Echo.
//
// Layout convention (Tiny RDM-style backend + DDD domains):
//
//	backend/services/      — Wails-bound API facades
//	backend/storage/       — SQLite, sqlc store, config.json
//	backend/types/         — shared Wails/JSON DTOs
//	backend/<domain>/      — domain code (watch, vdiupload, english/content, …)
//	cmd/<cli>/             — delivery adapters for side-car domains
//	cloud/<service>/       — cloud adapters when needed
package consts

const (
	// English is the English Learning domain: vocabulary, reading, SRS, Anki import.
	// Route prefix: "english.*". Code currently lives across content/, store/, appservice/.
	English = "english"

	// Watch is Echo Watch: local VDI client screen OCR + rule alerts + Cloudflare push.
	// Delivery: cmd/echo-watch, cloud/echo-watch. Code: backend/watch.
	Watch = "watch"

	// VDIUpload is automated host→guest file transfer for H3C Workspace (and later RDP/Citrix/Horizon).
	// Delivery: cmd/vdi-upload. Code: backend/vdiupload. Design: docs/domains/vdiupload/DESIGN.md.
	VDIUpload = "vdiupload"

	// Sudoku is the YouDoSudoku-backed puzzle game in the Agent window.
	Sudoku = "sudoku"

	// Obsidian is the Local REST API vault browser in the Management window.
	Obsidian = "obsidian"
)

// All lists known domain identifiers for discovery and tooling.
var All = []string{English, Watch, VDIUpload, Sudoku, Obsidian}

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
			CodeRoot:    "backend/english/content, backend/storage/store, backend/services/english",
			DocsRoot:    "docs/domains/english, docs/design-and-plan.md, docs/chat-learning-flows-plan.md",
			Description: "Vocabulary, articles, Anki import, SRS; first Wails Management/Agent domain.",
		},
		{
			ID:          Watch,
			Name:        "Echo Watch",
			Kind:        "sidecar",
			CodeRoot:    "backend/watch, cmd/echo-watch, cloud/echo-watch",
			DocsRoot:    "docs/domains/watch, cloud/echo-watch/README.md",
			Description: "Physical-PC screen watch of VDI client windows; OCR rules → Cloudflare → iOS PWA push.",
		},
		{
			ID:          VDIUpload,
			Name:        "VDI File Upload",
			Kind:        "sidecar",
			CodeRoot:    "backend/vdiupload, cmd/vdi-upload",
			DocsRoot:    "docs/domains/vdiupload",
			Description: "Automate file upload into H3C Workspace VDI via clipboard + Win32 UI automation.",
		},
		{
			ID:          Sudoku,
			Name:        "Sudoku",
			Kind:        "desktop-agent",
			CodeRoot:    "backend/sudoku, backend/services/sudoku.go, backend/storage",
			DocsRoot:    "docs/domains/sudoku",
			Description: "YouDoSudoku-backed puzzle sessions in the Agent window.",
		},
		{
			ID:          Obsidian,
			Name:        "Obsidian",
			Kind:        "desktop-agent",
			CodeRoot:    "backend/obsidian, backend/services/obsidian.go",
			DocsRoot:    "docs/domains/obsidian",
			Description: "Vault browse, note edit, and search via Obsidian Local REST API.",
		},
	}
}
