package instructions

import (
	"fmt"
	"sort"
	"time"

	"github.com/znwng/cdsl/internal/core"
)

// ========================================
// Reset Processing
// ========================================

// ProcessReset moves all configured components to their default positions.
//
// Components are processed in ascending order of their configured IDs. A
// two-second delay is inserted between each movement to allow the hardware
// to settle before the next command is sent.
func ProcessReset(config *core.Config, instruction core.Instruction) {
	components, ok := (*config)["component"].(map[string]any)
	if !ok {
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

	sort.Slice(componentList, func(i, j int) bool {
		return componentList[i].id < componentList[j].id
	})

	for _, component := range componentList {
		value := float32(component.defaultValue)
		valueText := fmt.Sprintf("%d", component.defaultValue)

		fmt.Printf(
			"%s = %d\n",
			component.name,
			component.defaultValue,
		)

		if !ExecuteMove(
			config,
			instruction,
			component.name,
			valueText,
			value,
		) {
			return
		}

		time.Sleep(2 * time.Second)
	}
}
