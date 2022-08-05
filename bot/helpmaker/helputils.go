package helpmaker

import (
	"sort"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

const (
	HelpPrifex = "help_"
)

func makeHelpBtn(modName string) gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{
		Text:         strings.Title(modName),
		CallbackData: HelpPrifex + modName,
	}
}

func getSortedModules() []string {
	mods := make([]string, 0)
	for module := range defaultHelper.ModHelp {
		mods = append(mods, module)
	}
	sort.Strings(mods)
	return mods
}
