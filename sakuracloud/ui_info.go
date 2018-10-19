package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/packer/packer"
)

var uiMsgPerStep = "%s step: %s %s"

func stepStartMsg(ui packer.Ui, debug bool, stepName string) {
	ui.Say(fmt.Sprintf(uiMsgPerStep, "-->", stepName, "start"))
}

func stepEndMsg(ui packer.Ui, debug bool, stepName string) {
	ui.Say(fmt.Sprintf(uiMsgPerStep, "<--", stepName, "end"))
}
