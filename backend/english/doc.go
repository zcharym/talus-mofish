// Package english documents the English Learning bounded context.
//
// Physical code lives in backend/english/content, backend/storage/store, and
// backend/services/english; this package anchors the domain ID for tooling.
//
// See docs/domains/english/README.md.
package english

import "github.com/songwei.ma/talus-mofish/backend/consts"

// ID is the canonical domain identifier.
const ID = consts.English
