package desktop

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDesktopItems(t *testing.T) {
	wd := os.Getenv("WORKSPACE_DIR")

	tests := []struct {
		name          string
		path          string
		expectedItems []Item
		expectedErr   error
	}{
		{
			name: "Success test",
			path: wd + "/test/data/app",
			expectedItems: []Item{
				{
					Shebang: "#!/usr/bin/env xdg-open",
					Path:    wd + "/test/data/app/share/applications/app.desktop",
					DesktopEntry: Section{
						Terminal:      false,
						StartupNotify: true,
						Name:          "MyApp",
						Exec:          "myapp %F",
						Type:          "Application",
						Icon:          "myapp-icon",
						Version:       "2.3",
						Comment:       "An example utility that performs useful tasks.",
						Categories:    "Utility;System;Development;",
						MimeType:      "application/json;application/xml;application/x-tar;",
						Actions:       "RunTask;Configure;",
						GenericName:   "Utility Application",
						TryExec:       "myapp",
						Keywords:      "utility;system;task;helper;",
					},
					Entries: map[string]Section{
						"Desktop Action Configure": {
							Name: "Configure the application",
							Exec: "myapp --configure",
						},
						"Desktop Action RunTask": {
							Name: "Run a predefined task",
							Exec: "myapp --run-task %F",
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := NewDesktop()
			items, err := d.GetDesktopItems(test.path)
			assert.ErrorIs(t, err, test.expectedErr)
			assert.Equal(t, test.expectedItems, items)
		})
	}
}
