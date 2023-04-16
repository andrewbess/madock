package setup

import (
	"github.com/faradey/madock/src/cli/fmtc"
	"github.com/faradey/madock/src/configs/projects"
	"github.com/faradey/madock/src/versions/pwa"
)

func PWA(projectName string, projectConfig map[string]string, continueSetup bool) {
	if continueSetup {
		toolsDefVersions := pwa.GetVersions()
		projects.SetEnvForProject(projectName, toolsDefVersions, projectConfig)
		fmtc.SuccessLn("\n" + "Finish set up environment")
	}
}
