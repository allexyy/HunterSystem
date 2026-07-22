package habit

import "strings"

type HabitData struct {
	Title       string
	Description string
	Difficult   string
	StatCodes   []string
}

func NewHabitData(message string) HabitData {
	habitText := strings.Split(message, "|")
	return HabitData{
		Title:       strings.TrimSpace(strings.TrimPrefix(habitText[0], "/newhabit")),
		Description: strings.TrimSpace(habitText[1]),
		Difficult:   strings.TrimSpace(habitText[2]),
		StatCodes:   habitText[3:],
	}
}
