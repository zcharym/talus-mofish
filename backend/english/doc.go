// Package english is the English Learning bounded context.
//
// Persistence lives in Repository; Anki import lives in content/; the Wails
// façade is backend/services/english.go.
//
// See docs/domains/english/README.md.
package english

import "github.com/songwei.ma/talus-mofish/backend/consts"

// ID is the canonical domain identifier.
const ID = consts.English
