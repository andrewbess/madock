package configs

import (
	"bytes"
	"encoding/xml"
	"github.com/faradey/madock/src/helper/paths"
	"github.com/go-xmlfmt/xmlfmt"
	"log"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"
)

func Save(file string, data map[string]string) {
	resultData := make(map[string]interface{})
	for key, value := range data {
		resultData[key] = value
	}
	resultMapData := SetXmlMap(resultData)
	w := &bytes.Buffer{}
	w.WriteString(xml.Header)
	err := MarshalXML(resultMapData, xml.NewEncoder(w), "scopes/default")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(file, []byte(xmlfmt.FormatXML(w.String(), "", "    ", true)), 0755)
	if err != nil {
		log.Fatalf("Unable to write file: %v", err)
	}
}

func IsHasConfig(projectName string) bool {
	PrepareDirsForProject(projectName)
	if paths.IsFileExist(paths.GetExecDirPath() + "/projects/" + projectName + "/config.xml") {
		return true
	}

	return false
}

func IsHasNotConfig() bool {
	if !paths.IsFileExist(paths.GetExecDirPath() + "/projects/" + GetProjectName() + "/config.xml") {
		return true
	}
	return false
}

func GeneralConfigMapping(mainConf map[string]string, targetConf map[string]string) {
	if len(mainConf) > 0 {
		for index, val := range mainConf {
			if v, ok := targetConf[index]; !ok || v == "" {
				targetConf[index] = val
			}
		}
	}
}

func ConfigMapping(mainConf map[string]string, targetConf map[string]string) {
	if len(targetConf) > 0 && len(mainConf) > 0 {
		for index, val := range mainConf {
			if _, ok := targetConf[index]; !ok {
				targetConf[index] = val
			}
		}
	}
}

func ReplaceConfigValue(str string) string {
	projectConf := GetCurrentProjectConfig()
	osArch := runtime.GOARCH
	arches := map[string]string{"arm64": "aarch64"}

	if arch, ok := arches[osArch]; ok {
		osArch = arch
	} else {
		osArch = "x86-64"
	}
	for key, val := range projectConf {
		str = strings.Replace(str, "{{{"+key+"}}}", val, -1)
	}

	str = strings.Replace(str, "{{{OSARCH}}}", osArch, -1)

	usr, err := user.Current()
	if err == nil {
		str = strings.Replace(str, "{{{UID}}}", usr.Uid, -1)
		str = strings.Replace(str, "{{{UNAME}}}", usr.Username, -1)
		str = strings.Replace(str, "{{{GUID}}}", usr.Gid, -1)
		gr, _ := user.LookupGroupId(usr.Gid)
		str = strings.Replace(str, "{{{UGROUP}}}", gr.Name, -1)
	} else {
		log.Fatal(err)
	}

	r := regexp.MustCompile("(?ism)<<<iftrue>>>(.*?)<<<endif>>>")
	str = r.ReplaceAllString(str, "$1")
	r = regexp.MustCompile("(?ism)<<<iffalse>>>.*?<<<endif>>>")
	str = r.ReplaceAllString(str, "")

	var onlyHosts []string

	hosts := projectConf["hosts"]
	if len(hosts) > 0 {
		for _, hostAndStore := range hosts {
			onlyHosts = append(onlyHosts, "- \""+strings.Split(hostAndStore, ":")[0]+":172.17.0.1\"")
		}
	}

	str = strings.Replace(str, "{{{HOST_GATEWAYS}}}", strings.Join(onlyHosts, "\n      "), -1)
	return str
}

func IsOption(name string) bool {
	for key := range GetCurrentProjectConfig() {
		if key == name {
			return true
		}
	}

	log.Fatalln("The option \"" + name + "\" doesn't exist.")

	return false
}
