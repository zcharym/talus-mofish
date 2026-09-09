package services

import (
	"github.com/songwei.ma/talus-mofish/backend/agent"
	"github.com/songwei.ma/talus-mofish/backend/types"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func init() {
	application.RegisterEvent[types.ConfigChanged]("config:changed")
	application.RegisterEvent[agent.StreamChunkEvent]("agent:stream-chunk")
	application.RegisterEvent[agent.TurnDoneEvent]("agent:turn-done")
	application.RegisterEvent[agent.TurnErrorEvent]("agent:turn-error")
	application.RegisterEvent[agent.TurnCancelledEvent]("agent:turn-cancelled")
}
