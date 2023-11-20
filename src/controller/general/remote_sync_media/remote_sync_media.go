package remote_sync_media

import (
	"github.com/faradey/madock/src/configs"
	"github.com/faradey/madock/src/ssh"
)

func Execute() {
	// TODO add CLI args
	projectConf := configs.GetCurrentProjectConfig()
	ssh.SyncMedia(projectConf["SSH_SITE_ROOT_PATH"])
}
