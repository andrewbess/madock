package stop

import (
	stopCustom "github.com/faradey/madock/src/controller/custom/stop"
	"github.com/faradey/madock/src/controller/general/proxy"
	stopMagento2 "github.com/faradey/madock/src/controller/magento/stop"
	stopPrestashop "github.com/faradey/madock/src/controller/prestashop/stop"
	stopPwa "github.com/faradey/madock/src/controller/pwa/stop"
	stopShopify "github.com/faradey/madock/src/controller/shopify/stop"
	stopShopware "github.com/faradey/madock/src/controller/shopware/stop"
	"github.com/faradey/madock/src/helper/configs"
	"github.com/faradey/madock/src/helper/paths"
)

func Execute() {
	projectConf := configs.GetCurrentProjectConfig()
	platform := projectConf["platform"]
	if platform == "magento2" {
		stopMagento2.Execute()
	} else if platform == "pwa" {
		stopPwa.Execute()
	} else if platform == "shopify" {
		stopShopify.Execute()
	} else if platform == "custom" {
		stopCustom.Execute()
	} else if platform == "shopware" {
		stopShopware.Execute()
	} else if platform == "prestashop" {
		stopPrestashop.Execute()
	}

	if len(paths.GetActiveProjects()) == 0 {
		proxy.Execute("stop")
	}
}
