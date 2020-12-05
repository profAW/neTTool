package main

import (
	log "github.com/sirupsen/logrus"
	"neTTool/helper"
	"neTTool/infrastructure"
	"neTTool/usecases"
	"os"
)

var config infrastructure.Configuration
var version = "1.0.5"

func main() {

	log.Info("### Welcome and remember 'never forget your towel' ###")
	log.Info("------------------------------------------------------")
	log.Info("neTTool-Version: " + version)
	log.Info("Check Licence")
	validLicence := helper.CheckLicence(helper.GetUser())

	if !validLicence {
		log.Fatal("Check Licence failed")
		log.Info("Press enter key to exit...")
		b := make([]byte, 1)
		_, _ = os.Stdin.Read(b)
		os.Exit(3)
	} else {
		log.Info("Valid Licence found")
		doAnalysis()
	}

}

func doAnalysis() {

	log.Info("------------------------------------------------------")
	log.Info("Load Configuration")

	ConfiSource := infrastructure.ConfigurationFromFS{}
	config = ConfiSource.LoadConfig()
	log.Info("Configuration loaded")
	log.Info("------------------------------------------------------")

	// TODO: Refactor in function with Unit-Test?
	if _, err := os.Stat("./results"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		errDir := os.Mkdir("./results", os.ModeDir)

		if errDir != nil {
			log.Error("Result-Folder could not be created, please run neTTool with admin permission.")
			log.Info("Press enter key to exit ...")
			b := make([]byte, 1)
			_, _ = os.Stdin.Read(b)
			os.Exit(3)
		}
	}

	log.Info("Start Analysis")
	log.Info("Get Network Data")

	data := usecases.UcGetNetworkData{}
	if config.Pcapfile != "" {
		data.Source = infrastructure.SavedPacketsAdapter{FileAndFolder: config.Pcapfile}
	} else {
		log.Error("No File in Config-File. Please provide file. ")
		log.Info("Press enter key to exit...")
		b := make([]byte, 1)
		_, _ = os.Stdin.Read(b)
		os.Exit(3)
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
	log.Info("    Finish PN Analysis")

	log.Info("Finish Analyse Network Connections")
	//fmt.Println("### Bye, and thank you for the fish ###")
	log.Info("Press Enter key to exit...")
	b := make([]byte, 1)
	_, _ = os.Stdin.Read(b)
}
