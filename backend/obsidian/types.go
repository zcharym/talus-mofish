package obsidian

// Status is the Local REST API GET / health response.
type Status struct {
	OK            string         `json:"ok"`
	Authenticated bool           `json:"authenticated"`
	Service       string         `json:"service"`
	Versions      StatusVersions `json:"versions"`
}

// StatusVersions reports plugin and Obsidian API versions.
type StatusVersions struct {
	Obsidian string `json:"obsidian"`
	Self     string `json:"self"`
}

// FileEntry is a vault file or directory from GET /vault/{dir}/.
type FileEntry struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
}

// Note is a markdown note from GET with Accept application/vnd.olrapi.note+json.
type Note struct {
	Path    string   `json:"path"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Stat    NoteStat `json:"stat"`
}

// NoteStat is filesystem metadata from the plugin.
type NoteStat struct {
	CTime float64 `json:"ctime"`
	MTime float64 `json:"mtime"`
	Size  float64 `json:"size"`
}

// SearchHit is one result from POST /search/simple/.
type SearchHit struct {
	Filename string        `json:"filename"`
	Score    float64       `json:"score"`
	Matches  []SearchMatch `json:"matches"`
}

// SearchMatch is a snippet around a search hit.
type SearchMatch struct {
	Context string `json:"context"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
}
