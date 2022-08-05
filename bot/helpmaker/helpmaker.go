package helpmaker

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Help struct {
	MainHelp         string
	ModHelp          map[string]string
	Parsemode        string
	PaginatedButtons map[int][][]gotgbot.InlineKeyboardButton
}

var defaultHelper = &Help{
	ModHelp:          make(map[string]string),
	PaginatedButtons: make(map[int][][]gotgbot.InlineKeyboardButton),
}

func SetMainHelp(help string, parsemode string) {
	defaultHelper.MainHelp = help
	defaultHelper.Parsemode = parsemode
}

func SetModuleHelp(module, help string) {
	defaultHelper.ModHelp[module] = help
}

func GetMainHelp() string {
	return defaultHelper.MainHelp
}

func GetParseMode() string {
	return defaultHelper.Parsemode
}

func GetModuleHelp(module string) string {
	return defaultHelper.ModHelp[module]
}

func MakeHelp() {
	defaultHelper.PaginatedButtons = PaginateModules(getSortedModules(), 3, 12, true)
}

func GetPageHelp(pageNumber int) [][]gotgbot.InlineKeyboardButton {
	return defaultHelper.PaginatedButtons[pageNumber]
}

func PaginateModules(modules []string, helpColumns int, maxNum int, omitEntry bool) map[int][][]gotgbot.InlineKeyboardButton {
	buttonMap := make(map[int][][]gotgbot.InlineKeyboardButton)
	rows := make([][]gotgbot.InlineKeyboardButton, 0)
	row := make([]gotgbot.InlineKeyboardButton, 0)
	var defEndBtn = gotgbot.InlineKeyboardButton{
		Text:         "Close",
		CallbackData: "help_close",
	}
	pageNum := 0
	for num, module := range modules {
		num += 1
		row = append(row, makeHelpBtn(module))
		if len(row) == helpColumns || num == (len(modules)) {
			rows = append(rows, row)
			row = make([]gotgbot.InlineKeyboardButton, 0)
		}
		if len(modules) > maxNum {
			if num >= maxNum && (num%maxNum == 0 || num == len(modules)) {
				pageNum += 1
				rows = append(rows, row)
				row = make([]gotgbot.InlineKeyboardButton, 0)
				if !omitEntry || pageNum != 1 {
					row = append(row, gotgbot.InlineKeyboardButton{
						Text:         "<",
						CallbackData: fmt.Sprintf("%sprev_%d", HelpPrifex, pageNum),
					})
				}
				row = append(row, defEndBtn)
				if !omitEntry || num != len(modules) {
					row = append(row, gotgbot.InlineKeyboardButton{
						Text:         ">",
						CallbackData: fmt.Sprintf("%snext_%d", HelpPrifex, pageNum),
					})
				}
				rows = append(rows, row)
				buttonMap[pageNum] = rows
				rows = make([][]gotgbot.InlineKeyboardButton, 0)
				row = make([]gotgbot.InlineKeyboardButton, 0)
			}
		} else if num == len(modules) {
			rows = append(rows, []gotgbot.InlineKeyboardButton{
				defEndBtn,
			})
			buttonMap[1] = rows
			break
		}
	}
	return buttonMap
}
