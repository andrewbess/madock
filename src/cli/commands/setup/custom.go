package setup

import (
	"fmt"
	"github.com/faradey/madock/src/cli/attr"
	"github.com/faradey/madock/src/cli/fmtc"
	"github.com/faradey/madock/src/configs/projects"
	"github.com/faradey/madock/src/docker/builder"
	"github.com/faradey/madock/src/paths"
	"github.com/faradey/madock/src/versions/custom"
)

func Custom(projectName string, projectConfig map[string]string, continueSetup bool) {
	toolsDefVersions := custom.GetVersions()

	if continueSetup {
		fmt.Println("")

		Php(&toolsDefVersions.Php)
		Db(&toolsDefVersions.Db)
		Composer(&toolsDefVersions.Composer)
		SearchEngine(&toolsDefVersions.SearchEngine)
		if toolsDefVersions.SearchEngine == "Elasticsearch" {
			Elastic(&toolsDefVersions.Elastic)
		} else {
			OpenSearch(&toolsDefVersions.OpenSearch)
		}

		Redis(&toolsDefVersions.Redis)
		RabbitMQ(&toolsDefVersions.RabbitMQ)
		Hosts(projectName, &toolsDefVersions.Hosts, projectConfig)

		projects.SetEnvForProject(projectName, toolsDefVersions, projectConfig)
		paths.MakeDirsByPath(paths.GetExecDirPath() + "/projects/" + projectName + "/backup/db")

		fmtc.SuccessLn("\n" + "Finish set up environment")
		fmtc.ToDoLn("Optionally, you can configure SSH access to the development server in order ")
		fmtc.ToDoLn("to synchronize the database and media files. Enter SSH data in ")
		fmtc.ToDoLn(paths.GetExecDirPath() + "/projects/" + projectName + "/env.txt")
	}

	builder.Down(attr.Options.WithVolumes)
	builder.StartCustom(attr.Options.WithChown, projectConfig)
}
