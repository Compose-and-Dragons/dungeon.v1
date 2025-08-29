package helpers

import (
	"fmt"

	"github.com/micro-agent/micro-agent-go/agent/msg"
	"github.com/micro-agent/micro-agent-go/agent/mu"
	"github.com/micro-agent/micro-agent-go/agent/ui"
)

func DisplayHistory(selectedAgent mu.Agent) {
	// remove the /debug part from the input
	fmt.Println()
	ui.Println(ui.Red, "📝 Messages history / Conversational memory:")
	for i, message := range selectedAgent.GetMessages() {
		printableMessage, err := msg.MessageToMap(message)
		if err != nil {
			ui.Printf(ui.Red, "Error converting message to map: %v\n", err)
			continue
		}
		ui.Print(ui.Cyan, "-", i, " ")
		ui.Print(ui.Orange, printableMessage["role"], ": ")
		ui.Println(ui.Cyan, printableMessage["content"])
	}
	ui.Println(ui.Red, "📝 End of the messages")
	fmt.Println()
}
