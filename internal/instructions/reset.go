package instructions

import (
	"fmt"
	"sort"
	"time"

	"github.com/kalexion/cdsl/internal/core"
	"github.com/kalexion/cdsl/internal/customTypes"
)

// ========================================
// Reset Processing
// ========================================

// ProcessReset moves all configured components to their default positions.
//
// Components are processed in ascending order of their configured IDs. A
// two-second delay is inserted between each movement to allow the hardware
// to settle before the next command is sent.
func ProcessReset(config *core.Config, instruction customTypes.Instruction) {
	components, ok := (*config)["component"].(map[string]any)
	if !ok {
		core.DiagnosticsError("No valid component configuration found", instruction)
		return
	}

	type componentInfo struct {
		name         string
		id           int64
		defaultValue int64
	}

	componentList := make([]componentInfo, 0, len(components))

	for name, rawComponent := range components {
		component, ok := rawComponent.(map[string]any)
		if !ok {
			continue
		}

		id, idOK := component["id"].(int64)
		defaultValue, defaultOK := component["default"].(int64)

		if !idOK || !defaultOK {
			continue
		}

		componentList = append(componentList, componentInfo{
			name:         name,
			id:           id,
			defaultValue: defaultValue,
		})
	}

	if len(componentList) == 0 {
		core.DiagnosticsError("No valid components available to reset", instruction)
		return
	}

	sort.Slice(componentList, func(i, j int) bool {
		return componentList[i].id < componentList[j].id
	})

	for i, component := range componentList {
		value := float32(component.defaultValue)
		valueText := fmt.Sprintf("%d", component.defaultValue)

		if !ExecuteMove(
			config,
			instruction,
			component.name,
			valueText,
			value,
		) {
			return
		}

		if i < len(componentList)-1 {
			time.Sleep(2 * time.Second)
		}
	}

	core.DiagnosticsSuccess("All components reset to their default positions", instruction)
}
