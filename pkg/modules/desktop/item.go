package desktop

import (
	"os"
	"strings"
)

// Item is a parsed desktop file
type Item struct {
	Shebang      string
	Path         string
	DesktopEntry Section
	Entries      map[string]Section
}

// NewItemFromPath returns desktop item from provided path
func NewItemFromPath(itemPath string) (Item, error) {
	data, err := os.ReadFile(itemPath) //nolint
	if err != nil {
		return Item{}, err
	}

	item := Item{Path: itemPath, Entries: make(map[string]Section)}
	item.parse(data)

	return item, nil
}

// NewItemFromData returns desktop item from provided binary data
func NewItemFromData(data []byte) Item {
	item := Item{Entries: make(map[string]Section)}
	item.parse(data)
	return item
}

func (item *Item) parse(data []byte) {
	splittedFile := strings.Split(string(data), "\n")
	sectionStart := -1

	for i, line := range splittedFile {
		if i == 0 && strings.HasPrefix(line, "#!") {
			item.Shebang = line
		}

		if strings.HasPrefix(line, "[") {
			if sectionStart == -1 {
				sectionStart = i
				continue
			}
			item.fillSection(splittedFile[sectionStart:i])
			sectionStart = i
		}
	}

	item.fillSection(splittedFile[sectionStart:])
}

func (item *Item) fillSection(section []string) {
	if len(section) == 0 {
		return
	}

	if strings.Contains(section[0], "[Desktop Entry]") {
		item.DesktopEntry.Parse(section)
		return
	}

	sectionTitle := strings.Trim(section[0], "[] \t\n\r")
	item.Entries[sectionTitle] = NewSection(section)
}
