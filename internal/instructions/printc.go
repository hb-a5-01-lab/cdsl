package instructions

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/znwng/cdsl/internal/core"
)

type componentInfo struct {
	name string
	id   int64
	min  int64
	max  int64
}

// ProcessPrintc prints the configuration in a tabular format
func ProcessPrintc(config *core.Config) {
	port, err := core.ArduinoPort(config)
	if err != nil {
		fmt.Printf("port = error: %v\n", err)
	} else {
		fmt.Printf("port = %q\n", port)
	}

	baudRate, err := core.ArduinoBaudRate(config)
	if err != nil {
		fmt.Printf("baud_rate = error: %v\n", err)
	} else {
		fmt.Printf("baud_rate = %d\n", baudRate)
	}

	fmt.Println()
	fmt.Println("Components:")

	components, ok := (*config)["component"].(map[string]any)
	if !ok || len(components) == 0 {
		fmt.Println("No components configured")
		return
	}

	componentList := make([]componentInfo, 0, len(components))

	for name, rawComponent := range components {
		component, ok := rawComponent.(map[string]any)
		if !ok {
			continue
		}

		id, idOK := component["id"].(int64)
		min, minOK := component["min"].(int64)
		max, maxOK := component["max"].(int64)

		if !idOK || !minOK || !maxOK {
			continue
		}

		componentList = append(componentList, componentInfo{
			name: name,
			id:   id,
			min:  min,
			max:  max,
		})
	}

	sort.Slice(componentList, func(i, j int) bool {
		return componentList[i].id < componentList[j].id
	})

	idWidth := len("c_id")
	nameWidth := len("c_name")
	rangeWidth := len("c_range")

	for _, component := range componentList {
		idWidth = max(idWidth, len(strconv.FormatInt(component.id, 10)))

		rangeValue := fmt.Sprintf("%d-%d", component.min, component.max)

		nameWidth = max(nameWidth, len(component.name))
		rangeWidth = max(rangeWidth, len(rangeValue))
	}

	fmt.Printf(
		"%-*s | %-*s | %-*s\n",
		idWidth,
		"c_id",
		nameWidth,
		"c_name",
		rangeWidth,
		"c_range",
	)

	tableWidth := idWidth + 3 + nameWidth + 3 + rangeWidth
	fmt.Println(strings.Repeat("-", tableWidth))

	for _, component := range componentList {
		rangeValue := fmt.Sprintf("%d-%d", component.min, component.max)

		fmt.Printf(
			"%-*d | %-*s | %-*s\n",
			idWidth,
			component.id,
			nameWidth,
			component.name,
			rangeWidth,
			rangeValue,
		)
	}
}
