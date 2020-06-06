package main

import (
	log "github.com/sirupsen/logrus"
	"neTTool/infrastructure"
	"neTTool/usecases"
	"os"
	"os/user"
	"strings"
)

var config infrastructure.Configuration
var version = "1.0.2"

func checkLicence() bool {

	localUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Current User
	//fmt.Println("Username: " + localUser.Username)

	validLicence := false
	//fmt.Println(strings.Contains(localUser.Username, "wenzela"))
	if strings.Contains(localUser.Username, "wenzela") || strings.Contains(localUser.Username, "atlabor") {
		validLicence = true
	}

	return validLicence
}

func main() {

	log.Info("### Welcome and remember 'never forget your towel' ###")
	log.Info("------------------------------------------------------")
	log.Info("neTTool-Version: " + version)
	log.Info("Check Licence")
	validLicence := checkLicence()

	if !validLicence {
		log.Fatal("Check Licence falied")
		log.Info("Press any key to exit...")
		b := make([]byte, 1)
		os.Stdin.Read(b)
		os.Exit(3)
	} else {
		log.Info("Valid Licence found")
	}
	log.Info("------------------------------------------------------")
	log.Info("Load Configuraiton")

	ConfiSource := infrastructure.ConfigurationFromFS{}
	config = ConfiSource.LoadConfig()
	log.Info("Configuraiton loaded")
	log.Info("------------------------------------------------------")

	os.Mkdir("./results", os.ModeDir)

	log.Info("Start Analysis")
	log.Info("Get Netzwork Data")

	sourceA := infrastructure.SavedPacketsAdapter{FileAndFolder: config.Pcapfile}

	mydevice := config.InterfaceID
	sourceB := infrastructure.LivePacketsAdapter{Device: mydevice, SnapshotLen: config.SnapshotLen, Promiscuous: config.Promiscuous, Timeout: config.Timeout}

	data := usecases.UcGetNetworkData{}
	if config.StoredData {
		data.Source = sourceA
	} else {
		data.Source = sourceB
	}
	packetSource := data.Read()

	log.Info("    Start Analyse Network Connections")
	graphDestination := infrastructure.SaveConnectionGraphToFsAdapter{FileAndFolder: "./results/networkgraph.gv"}
	analysis := usecases.UcConnectionAnalysis{Destination: graphDestination}
	connection := analysis.CreateConnectionList(packetSource)
	connectionGraph := analysis.MakeConnetionGraph(connection)
	analysis.ExportConnectionGraph(connectionGraph)
	log.Info("    Finish Analyse Network Connections")

	log.Info("    Start PN Analysis")

	analysisPN := usecases.UcProfiNETAnalysis{}
	pnData := analysisPN.GetProfiNetData(packetSource)
	pnData = analysisPN.CalcProfiNetDeltaTimeInMS(pnData)
	pnResultExport := infrastructure.SavePNGraphToFsAdapter{FileAndFolder: ""}
	pnResultExport.PlotData(pnData)
	log.Info("        PN-Analysis created")
	log.Info("    Finsh PN Analysis")

	log.Info("Finish Analyse Network Connections")
	//fmt.Println("### Bye, and thank you for the fish ###")
	log.Info("Press Enter key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
