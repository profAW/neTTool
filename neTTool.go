package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"neTTool/domain"
	"neTTool/helper"
	"neTTool/infrastructure"
	"neTTool/usecases"
	"os"
	"time"
)

type fileSettings struct {
	pcapSource        string
	pcapFileName      string
	resultDestination string
}

var version = "neTTool-Version: 1.3.0"
var mySettings fileSettings

func main() {

	fmt.Println("╱╱╱╱╱╱╱╱╱╱╭━━━━┳━━━━╮╱╱╱╱╭╮╱╱╱╱╱")
	fmt.Println("╱╱╱╱╱╱╱╱╱╱┃╭╮╭╮┃╭╮╭╮┃╱╱╱╱┃┃╱╱╱╱╱")
	fmt.Println("╱╱╱╱╱╭━╮╭━┻┫┃┃╰┻╯┃┃┣┻━┳━━┫┃╱╱╱╱╱")
	fmt.Println("╱╱╱╱╱┃╭╮┫┃━┫┃┃╱╱╱┃┃┃╭╮┃╭╮┃┃╱╱╱╱╱")
	fmt.Println("╱╱╱╱╱┃┃┃┃┃━┫┃┃╱╱╱┃┃┃╰╯┃╰╯┃╰╮╱╱╱╱╱")
	fmt.Println("╱╱╱╱╱╰╯╰┻━━╯╰╯╱╱╱╰╯╰━━┻━━┻━╯╱╱╱╱╱")
	fmt.Println("----------------------------------")
	fmt.Println(version)

	myApp := app.New()
	myWindow := myApp.NewWindow("neTTool")
	myWindow.Resize(fyne.NewSize(800, 400))
	headerLabel := canvas.Text{Text: "neTTool", TextSize: 40.0, Alignment: fyne.TextAlignCenter}

	subHeaderLabel0 := canvas.Text{Text: "A small tool to analysis profiNET connections", TextSize: 13.0, Alignment: fyne.TextAlignCenter}
	subHeaderLabel1 := canvas.Text{Text: version, TextSize: 8.0, Alignment: fyne.TextAlignCenter}
	subHeaderLabel2 := canvas.Text{Text: "HAW Hamburg", TextSize: 10.0, Alignment: fyne.TextAlignTrailing}
	subHeaderLabel3 := canvas.Text{Text: "Prof. Dr.-Ing. A. Wenzel", TextSize: 10.0, Alignment: fyne.TextAlignTrailing}

	path := binding.NewString()
	sourcePathLabel := widget.NewLabelWithData(path)

	pathDestination := binding.NewString()
	destinationPathLabel := widget.NewLabelWithData(pathDestination)

	analysisStatus := binding.NewString()
	analysisStatusLabel := widget.NewLabelWithData(analysisStatus)

	loadFileButtonText := "Step 1 - Select a pcap-file for analysis"
	doAnalysisButtonText := "Step 2 - Perform analysis                        "
	saveBuottonText := "Step 3 - Select folder to save results    "

	loadAnalysisFileButton := widget.NewButtonWithIcon(loadFileButtonText, theme.DocumentIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			path.Set("Source-File loaded: " + reader.URI().Path())
			mySettings.pcapSource = reader.URI().Path()
			mySettings.pcapFileName = reader.URI().Name()
			helper.RemoveContents("./results")
			analysisStatus.Set("Analysis open")
		}, myWindow)
	})

	doAnalysisButton := widget.NewButtonWithIcon(doAnalysisButtonText, theme.MediaPlayIcon(), func() {
		analysisStatus.Set("Analysis in progress")
		doAnalysis()
		analysisStatus.Set("Analysis done")
	})

	saveResultsButton := widget.NewButtonWithIcon(saveBuottonText, theme.DocumentSaveIcon(), func() {
		dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
			if err == nil && folder == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			currentTime := time.Now()
			mySettings.resultDestination = folder.Path() + "/neTTool_" + mySettings.pcapFileName + "_" + currentTime.Format("2006_01_02_15_04_05") + ".zip"
			pathDestination.Set("Destination-File saved @: " + mySettings.resultDestination)
			infrastructure.ZipResults("./results", mySettings.resultDestination)
		}, myWindow)
	})

	loadAnalysisFileButton.Alignment = widget.ButtonAlignLeading
	doAnalysisButton.Alignment = widget.ButtonAlignLeading
	saveResultsButton.Alignment = widget.ButtonAlignLeading

	centered := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), &headerLabel, layout.NewSpacer())

	myWindow.SetContent(container.New(layout.NewVBoxLayout(), centered, &subHeaderLabel0, &subHeaderLabel1, &subHeaderLabel2, &subHeaderLabel3, layout.NewSpacer(), loadAnalysisFileButton, sourcePathLabel, doAnalysisButton, analysisStatusLabel, saveResultsButton, destinationPathLabel, layout.NewSpacer()))

	myWindow.Show()
	myApp.Run()

}

func checkIfPnIsPresent(Data map[string]domain.CommonConnection) bool {
	// Test ob es PN-Daten gibt
	/*
			Iteration über connectionList
			Test ob EtherTyp PN enthalen ist
		    falls ja, boolvar auf true setzen
	*/

	var doPnAnaylsis = false
	for _, element := range Data {
		if element.EthernetType == "8892" {
			doPnAnaylsis = true
		}
	}
	return doPnAnaylsis
}

func doAnalysis() {

	if !helper.Exists("./results") {
		// path/to/whatever does not exist
		errDir := os.Mkdir("./results", os.ModeDir)

		if errDir != nil {
			fmt.Println("Result-Folder could not be created, please run neTTool with admin permission.")
			fmt.Println("Press enter key to exit ...")
			helper.CloseApplicationWithError()
		}
	}

	data := usecases.UcGetNetworkData{}
	data.Source = infrastructure.SavedPacketsAdapter{FileAndFolder: mySettings.pcapSource}
	packetSource := data.Read()
	connectionsList := data.CreateNetworkData(packetSource)

	graphDestination := infrastructure.SaveConnectionGraphToFsAdapter{FileAndFolder: "./results/networkgraph.gv"}
	analysis := usecases.UcConnectionAnalysis{Destination: graphDestination}
	connectionGraph := analysis.MakeConnetionGraph(connectionsList)

	analysis.ExportConnectionGraph(connectionGraph)

	if checkIfPnIsPresent(connectionsList) {
		analysisPN := usecases.UcProfiNETAnalysis{}
		connectionsList = analysisPN.CalcProfiNetDeltaTimeInMS(connectionsList)
		pnResultExport := infrastructure.SavePNGraphToFsAdapter{FileAndFolder: ""}
		pnResultExport.PlotData(connectionsList)
	}

}
