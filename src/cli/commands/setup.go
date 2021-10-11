package commands

import (
	"bufio"
	"fmt"
	"github.com/faradey/madock/src/cli/fmtc"
	"github.com/faradey/madock/src/configs"
	"github.com/faradey/madock/src/docker/builder"
	"github.com/faradey/madock/src/paths"
	"github.com/faradey/madock/src/versions"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func Setup() {
	configs.IsHasConfig()
	builder.DownAll()
	fmtc.SuccessLn("Start set up environment")
	toolsDefVersions := versions.GetVersions()

	setupPhp(&toolsDefVersions.Php)
	setupDB(&toolsDefVersions.Db)
	setupComposer(&toolsDefVersions.Composer)
	setupElastic(&toolsDefVersions.Elastic)
	setupRedis(&toolsDefVersions.Redis)
	setupRabbitMQ(&toolsDefVersions.RabbitMQ)

	configs.SetEnvForProject(toolsDefVersions)

	createProjectNginxConf()
	createProjectNginxDockerfile()
	paths.MakeDirsByPath(paths.GetExecDirPath() + "/projects/" + paths.GetRunDirName() + "/backup/db")
	paths.MakeDirsByPath(paths.GetExecDirPath() + "/aruntime/projects/" + paths.GetRunDirName() + "/data/mysql")

	fmtc.SuccessLn("Finish set up environment")
}

func setupPhp(defVersion *string) {
	setTitleAndRecommended("PHP", defVersion)

	availableVersions := []string{"8.1", "8.0", "7.4", "7.3", "7.2", "7.1", "7.0"}

	for index, ver := range availableVersions {
		fmt.Println(strconv.Itoa(index+1) + ") " + ver)
	}

	invitation(defVersion)

	waiter(defVersion, availableVersions)
}

func setupDB(defVersion *string) {
	setTitleAndRecommended("MariaDB", defVersion)

	availableVersions := []string{"10.4", "10.3", "10.2", "10.1", "10.0"}

	for index, ver := range availableVersions {
		fmt.Println(strconv.Itoa(index+1) + ") " + ver)
	}

	invitation(defVersion)

	waiter(defVersion, availableVersions)
}

func setupComposer(defVersion *string) {
	setTitleAndRecommended("Composer", defVersion)

	availableVersions := []string{"1", "2"}

	for index, ver := range availableVersions {
		fmt.Println(strconv.Itoa(index+1) + ") " + ver)
	}
	invitation(defVersion)

	waiter(defVersion, availableVersions)
}

func setupElastic(defVersion *string) {
	setTitleAndRecommended("Elasticsearch", defVersion)

	availableVersions := []string{"7.10", "7.9", "7.7", "7.6", "6.8", "5.1"}

	for index, ver := range availableVersions {
		fmt.Println(strconv.Itoa(index+1) + ") " + ver)
	}

	invitation(defVersion)

	waiter(defVersion, availableVersions)
}

func setupRedis(defVersion *string) {
	setTitleAndRecommended("Redis", defVersion)

	availableVersions := []string{"6.0", "5.0"}

	for index, ver := range availableVersions {
		fmt.Println(strconv.Itoa(index+1) + ") " + ver)
	}

	invitation(defVersion)

	waiter(defVersion, availableVersions)
}

func setupRabbitMQ(defVersion *string) {
	setTitleAndRecommended("RabbitMQ", defVersion)
	availableVersions := []string{"3.8", "3.7"}
	prepareVersions(availableVersions)
	invitation(defVersion)
	waiter(defVersion, availableVersions)
}

func prepareVersions(availableVersions []string) {
	for index, ver := range availableVersions {
		fmt.Println(strconv.Itoa(index+1) + ") " + ver)
	}
}

func setTitleAndRecommended(title string, recommended *string) {
	fmtc.TitleLn(title)
	if *recommended != "" {
		fmt.Println("Recommended version: " + *recommended)
	}
}

func invitation(ver *string) {
	if *ver != "" {
		fmt.Println("Enter the item number or press Enter to select the recommended version")
	} else {
		fmt.Println("Enter the item number")
	}

	fmt.Print("> ")
}

func waiter(defVersion *string, availableVersions []string) {
	buf := bufio.NewReader(os.Stdin)
	sentence, err := buf.ReadBytes('\n')
	selected := strings.TrimSpace(string(sentence))
	if err != nil {
		log.Fatalln(err)
	} else {
		selectedInt, err := strconv.Atoi(selected)
		if selected == "" && *defVersion != "" {
			fmt.Println("Your choice: " + *defVersion)
		} else if selected != "" && err == nil && len(availableVersions) >= selectedInt {
			*defVersion = availableVersions[selectedInt-1]
			fmt.Println("Your choice: " + *defVersion)
		} else {
			fmtc.WarningLn("Choose one of the options offered")
			setupElastic(defVersion)
		}
	}
}

func createProjectNginxConf() {
	projectName := paths.GetRunDirName()
	nginxDefFile := paths.GetExecDirPath() + "/docker/nginx/conf/default.conf"
	b, err := os.ReadFile(nginxDefFile)
	if err != nil {
		log.Fatal(err)
	}

	paths.MakeDirsByPath(paths.GetExecDirPath() + "/projects/" + projectName + "/docker/nginx/conf")
	nginxFile := paths.GetExecDirPath() + "/projects/" + projectName + "/docker/nginx/conf/default.conf"
	err = ioutil.WriteFile(nginxFile, b, 0755)
	if err != nil {
		log.Fatalf("Unable to write file: %v", err)
	}
}

func createProjectNginxDockerfile() {
	projectName := paths.GetRunDirName()
	nginxDefFile := paths.GetExecDirPath() + "/docker/nginx/Dockerfile"
	b, err := os.ReadFile(nginxDefFile)
	if err != nil {
		log.Fatal(err)
	}

	paths.MakeDirsByPath(paths.GetExecDirPath() + "/projects/" + projectName + "/docker/nginx")
	err = ioutil.WriteFile(paths.GetExecDirPath()+"/projects/"+projectName+"/docker/nginx/Dockerfile", b, 0755)
	if err != nil {
		log.Fatalf("Unable to write file: %v", err)
	}
}
