package db

import (
	"encoding/json"
	"fmt"
	"github.com/faradey/madock/src/controller/general/remote_sync"
	"github.com/faradey/madock/src/helper/cli/attr"
	"github.com/faradey/madock/src/helper/cli/fmtc"
	"github.com/faradey/madock/src/helper/configs"
	"github.com/faradey/madock/src/helper/paths"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/sftp"
	"log"
	"os"
	"strings"
	"time"
)

type ArgsStruct struct {
	attr.Arguments
	Name        string   `long:"name" short:"n" description:"Name of the archive file"`
	IgnoreTable []string `long:"ignore-table" short:"i" description:"Ignore db table"`
}

func Execute() {
	args := getArgs()

	projectConf := configs.GetCurrentProjectConfig()
	conn := remote_sync.Connect(projectConf["SSH_AUTH_TYPE"], projectConf["SSH_KEY_PATH"], projectConf["SSH_PASSWORD"], projectConf["SSH_HOST"], projectConf["SSH_PORT"], projectConf["SSH_USERNAME"])

	remoteDir := projectConf["SSH_SITE_ROOT_PATH"]
	name := args.Name
	//TODO add options --db-user --db-password --db-name --db-host --db-port
	defer conn.Close()
	fmt.Println("")
	fmt.Println("Dumping and downloading DB is started")
	result := remote_sync.RunCommand(conn, "php -r \"\\$r1 = include('"+remoteDir+"/app/etc/env.php'); echo json_encode(\\$r1[\\\"db\\\"][\\\"connection\\\"][\\\"default\\\"]);\"")
	nOpenBrace := strings.Index(result, "{")
	if nOpenBrace != -1 {
		result = result[nOpenBrace:]
	} else {
		fmt.Println(result)
		log.Fatal("Failed to get database authentication data")
	}
	if len(result) > 2 {
		dbAuthData := remote_sync.RemoteDbStruct{}
		err := json.Unmarshal([]byte(result), &dbAuthData)
		if err != nil {
			fmt.Println(err)
		}
		curDateTime := time.Now().Format("2006-01-02_15-04-05")
		name = strings.TrimSpace(name)
		if len(name) > 0 {
			name += "_"
		}
		dumpName := "remote_" + name + curDateTime + ".sql.gz"

		ignoreTablesStr := ""
		ignoreTables := args.IgnoreTable
		if len(ignoreTables) > 0 {
			ignoreTablesStr = " --ignore-table=" + dbAuthData.Dbname + "." + strings.Join(ignoreTables, " --ignore-table="+dbAuthData.Dbname+".")
		}

		result = remote_sync.RunCommand(conn, "mysqldump -u \""+dbAuthData.Username+"\" -p\""+dbAuthData.Password+"\" -h "+dbAuthData.Host+" --quick --lock-tables=false --no-tablespaces --triggers"+ignoreTablesStr+" "+dbAuthData.Dbname+" | sed -e 's/DEFINER[ ]*=[ ]*[^*]*\\*/\\*/' | gzip > "+"/tmp/"+dumpName)
		sc, err := sftp.NewClient(conn)
		if err != nil {
			log.Fatal(err)
		}
		defer sc.Close()
		execPath := paths.GetExecDirPath()
		projectName := configs.GetProjectName()
		err = remote_sync.DownloadFile(sc, "/tmp/"+dumpName, execPath+"/projects/"+projectName+"/backup/db/"+dumpName, false, false)
		if err != nil {
			log.Fatal(err)
		}
		result = remote_sync.RunCommand(conn, "rm "+"/tmp/"+dumpName)
		fmt.Println("")
		fmtc.SuccessLn("A database dump was created and saved locally. To import a database dump locally run the command `madock db:import`")
	} else {
		fmt.Println("Failed to get database authentication data")
	}
}

func getArgs() *ArgsStruct {
	args := new(ArgsStruct)
	if len(os.Args) > 2 {
		argsOrigin := os.Args[2:]
		var err error
		_, err = flags.ParseArgs(args, argsOrigin)

		if err != nil {
			log.Fatal(err)
		}
	}

	return args
}
