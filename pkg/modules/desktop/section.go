package desktop

import (
	"fmt"
	"strconv"
	"strings"
)

// Section is a desktop file section
type Section struct {
	Terminal, StartupNotify, NoDisplay, Hidden   bool
	Name, Exec, Type, Icon                       string // required
	Version, Comment, Categories, MimeType, Path string
	Actions, GenericName, TryExec, Keywords      string
	XGnomeSettingsPanel, XKdeSubstituteUID       string
	XAppInstallPackage                           string
	XWindowIcon, XDesktopFileInstallVersion      string
	XSchemeHandlerMailto, XDBusService           string
}

// NewSection returns new section
func NewSection(lines []string) Section {
	s := Section{}
	s.Parse(lines)
	return s
}

// Parse returns parsed section
func (s *Section) Parse(lines []string) {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		splittedLine := strings.SplitN(line, "=", 2)
		if len(splittedLine) != 2 {
			continue
		}
		key, value := splittedLine[0], splittedLine[1]

		switch key {
		case "Name":
			s.Name = value
		case "GenericName":
			s.GenericName = value
		case "Type":
			s.Type = value
		case "Version":
			s.Version = value
		case "Comment":
			s.Comment = value
		case "Exec":
			s.Exec = value
		case "TryExec":
			s.TryExec = value
		case "Icon":
			s.Icon = value
		case "MimeType":
			s.MimeType = value
		case "Keywords":
			s.Keywords = value
		case "Categories":
			s.Categories = value
		case "Terminal":
			s.Terminal, _ = strconv.ParseBool(value)
		case "StartupNotify":
			s.StartupNotify, _ = strconv.ParseBool(value)
		case "NoDisplay":
			s.NoDisplay, _ = strconv.ParseBool(value)
		case "Hidden":
			s.Hidden, _ = strconv.ParseBool(value)
		case "Path":
			s.Path = value
		case "Actions":
			s.Actions = value
		case "X-GNOME-Settings-Panel":
			s.XGnomeSettingsPanel = value
		case "X-KDE-SubstituteUID":
			s.XKdeSubstituteUID = value
		case "X-AppInstall-Package":
			s.XAppInstallPackage = value
		case "X-Window-Icon":
			s.XWindowIcon = value
		case "X-Desktop-File-Install-Version":
			s.XDesktopFileInstallVersion = value
		case "X-SchemeHandler-mailto":
			s.XSchemeHandlerMailto = value
		case "X-DBus-Service":
			s.XDBusService = value
		default:
			fmt.Printf("Unknown key: %s\n", key)
		}
	}
}
