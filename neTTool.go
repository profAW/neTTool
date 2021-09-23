package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"neTTool/helper"
	"neTTool/infrastructure"
	"neTTool/usecases"
	"os"
)

var config infrastructure.Configuration
var version = "1.0.8"

func main() {

	fmt.Println("╱╱╱╱╱╱╱╱╱╱╭━━━━┳━━━━╮╱╱╱╱╭╮╱╱╱╱╱")
	fmt.Println("╱╱╱╱╱╱╱╱╱╱┃╭╮╭╮┃╭╮╭╮┃╱╱╱╱┃┃╱╱╱╱╱")
	fmt.Println("╱╱╱╱╱╭━╮╭━┻┫┃┃╰┻╯┃┃┣┻━┳━━┫┃╱╱╱╱╱")
	fmt.Println("╱╱╱╱╱┃╭╮┫┃━┫┃┃╱╱╱┃┃┃╭╮┃╭╮┃┃╱╱╱╱╱")
	fmt.Println("╱╱╱╱╱┃┃┃┃┃━┫┃┃╱╱╱┃┃┃╰╯┃╰╯┃╰╮╱╱╱╱╱")
	fmt.Println("╱╱╱╱╱╰╯╰┻━━╯╰╯╱╱╱╰╯╰━━┻━━┻━╯╱╱╱╱╱")
	fmt.Println("----------------------------------")
	fmt.Println("neTTool-Version: " + version)

	doAnalysis()

	helper.CloseApplicationWithOutError()

}

func doAnalysis() {

	fmt.Println("------------------------------------------------------")
	fmt.Println("Load Configuration")

	ConfiSource := infrastructure.ConfigurationFromFS{}
	config = ConfiSource.LoadConfig()

	if !helper.Exists(config.Pcapfile) {
		log.Error("Can not access the pcap-file from Configuration. Please check path and file.")
		log.Error("Configuration-File-Path-Name: ", config.Pcapfile)
		fmt.Println("Press enter key to exit...")
		helper.CloseApplicationWithError()
	}

	fmt.Println("Configuration loaded")
	fmt.Println("------------------------------------------------------")

	if !helper.Exists("./results") {
		// path/to/whatever does not exist
		errDir := os.Mkdir("./results", os.ModeDir)

		if errDir != nil {
			log.Error("Result-Folder could not be created, please run neTTool with admin permission.")
			fmt.Println("Press enter key to exit ...")
			helper.CloseApplicationWithError()
		}
	}

	fmt.Println("Start Analysis")
	fmt.Println("Get Network Data")

	data := usecases.UcGetNetworkData{}
	data.Source = infrastructure.SavedPacketsAdapter{FileAndFolder: config.Pcapfile}
	packetSource := data.Read()
	connectionsList := data.CreateNetworkData(packetSource)

	fmt.Println("    Start Analyse Network Connections")
	graphDestination := infrastructure.SaveConnectionGraphToFsAdapter{FileAndFolder: "./results/networkgraph.gv"}
	analysis := usecases.UcConnectionAnalysis{Destination: graphDestination}
	connectionGraph := analysis.MakeConnetionGraph(connectionsList)

	analysis.ExportConnectionGraph(connectionGraph)
	fmt.Println("    Finish Analyse Network Connections")

	fmt.Println("    Start PN Analysis")

	analysisPN := usecases.UcProfiNETAnalysis{}
	connectionsList = analysisPN.CalcProfiNetDeltaTimeInMS(connectionsList)
	pnResultExport := infrastructure.SavePNGraphToFsAdapter{FileAndFolder: ""}
	pnResultExport.PlotData(connectionsList)
	fmt.Println("        PN-Analysis created")
	fmt.Println("    Finish PN Analysis")

	fmt.Println("Finish Analyse Network Connections")
	fmt.Println("Press enter key to exit...")
}
