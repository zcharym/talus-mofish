// Package english documents the English Learning bounded context.
//
// Physical code still lives in internal/content, internal/store, and
// internal/appservice; this package anchors the domain ID for tooling.
//
// See docs/domains/english/README.md.
package english

import "github.com/songwei.ma/talus-mofish/internal/domain"

// ID is the canonical domain identifier.
const ID = domain.English
