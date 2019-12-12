package main

import (
	"fmt"
	"neTTool/infrastructure"
	"neTTool/usecases"
	"os"
	"os/user"
	"strings"
)

var config infrastructure.Configuration
var version = "1.0.1"

func checkLicence() bool {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Current User
	//fmt.Println("Username: " + user.Username)

	validLicence := false
	//fmt.Println(strings.Contains(user.Username, "wenzela"))
	if strings.Contains(user.Username, "wenzela") || strings.Contains(user.Username, "atlabor") {
		validLicence = true
	}

	return validLicence
}

func main() {

	fmt.Println("### Welcome and remember 'never forget your towel' ###")
	fmt.Println("------------------------------------------------------")
	fmt.Println("neTTool-Version: " + version)
	fmt.Println("Check Licence")
	validLicence := checkLicence()

	if !validLicence {
		fmt.Println("Check Licence falied")
		fmt.Printf("Press any key to exit...")
		b := make([]byte, 1)
		os.Stdin.Read(b)
		os.Exit(3)
	} else {
		fmt.Println("Valid Licence found")
	}

	fmt.Println("------------------------------------------------------")
	fmt.Println("Load Configuraiton")

	ConfiSource := infrastructure.ConfigurationFromFS{}
	config = ConfiSource.LoadConfig()
	fmt.Println("Configuraiton loaded")
	fmt.Println("------------------------------------------------------")

	os.Mkdir("./results", os.ModeDir)

	fmt.Println("Start Analysis")
	fmt.Println("Get Netzwork Data")

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

	fmt.Println("    Start Analyse Network Connections")
	graphDestination := infrastructure.SaveConnectionGraphToFsAdapter{FileAndFolder: "./results/networkgraph.gv"}
	analysis := usecases.UcConnectionAnalysis{Destination: graphDestination}
	conncetion := analysis.CreateConnectionList(packetSource)
	connectionGraph := analysis.MakeConnetionGraph(conncetion)
	analysis.ExportConnectionGraph(connectionGraph)
	fmt.Println("    Finish Analyse Network Connections")

	fmt.Println("    Start PN Analysis")

	analysisPN := usecases.UcProfiNETAnalysis{}
	pnData := analysisPN.GetProfiNetData(packetSource)
	pnData = analysisPN.CalcProfiNetDeltaTimeInMS(pnData)
	pnResultExport := infrastructure.SavePNGraphToFsAdapter{FileAndFolder: ""}
	pnResultExport.PlotData(pnData)
	fmt.Println("        PN-Analysis created")
	fmt.Println("    Finsh PN Analysis")

	fmt.Println("Finish Analyse Network Connections")
	fmt.Println("### Bye, and thank you for the fish ###")
	fmt.Printf("Press any key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
