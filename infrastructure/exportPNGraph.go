package infrastructure

import (
	"fmt"
	"neTTool/domain"
	"os"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// SavePNGraphToFsAdapter stores PN Analysis-Data to FS
type SavePNGraphToFsAdapter struct {
	FileAndFolder string
}

// PlotData exports the PN Analysis Result to FS
func (e SavePNGraphToFsAdapter) PlotData(Daten map[string]domain.ProfinetData) {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "PN Cycle Time from Source to Destination"
	p.Y.Label.Text = "Time [ms]"

	// Make boxes for our data and add them to the plot.
	statitics := []string{}
	xtext := make(map[int]string)
	loc := 0.0
	i := 0
	for k, con := range Daten {

		if len(con.DeltaTS) > 0 {

			plotData := make(plotter.Values, len(con.DeltaTS))
			for i, s := range con.DeltaTS {
				plotData[i] = s
			}
			w := vg.Points(20)
			b0, err := plotter.NewBoxPlot(w, loc, plotData)
			if err != nil {
				panic(err)
			}
			statitics = append(statitics, "----------------------------------------------")

			statitics = append(statitics, string(k))
			statitics = append(statitics, "Max            : "+fmt.Sprintf("%2f", b0.Max)+"ms")
			statitics = append(statitics, "Oberer Wisker  : "+fmt.Sprintf("%2f", b0.AdjHigh)+"ms")
			statitics = append(statitics, "75% Quantil    : "+fmt.Sprintf("%2f", b0.Quartile3)+"ms")
			statitics = append(statitics, "Median         : "+fmt.Sprintf("%2f", b0.Median)+"ms")
			statitics = append(statitics, "25% Quantil    : "+fmt.Sprintf("%2f", b0.Quartile1)+"ms")
			statitics = append(statitics, "Unterer Wisker : "+fmt.Sprintf("%2f", b0.AdjLow)+"ms")
			statitics = append(statitics, "Min            : "+fmt.Sprintf("%2f", b0.Min)+"ms")

			p.Add(b0)
			xtext[i] = strings.ReplaceAll(k, "->", " \n -->\n")
			loc = loc + 1
			i = i + 1
		}

	}

	p.Y.Min = -0.1
	//p.Y.Max = 6

	switch i {
	case 1:
		p.NominalX(xtext[0])
	case 2:
		p.NominalX(xtext[0], xtext[1])
	case 3:
		p.NominalX(xtext[0], xtext[1], xtext[2])
	case 4:
		p.NominalX(xtext[0], xtext[1], xtext[2], xtext[3])
	case 5:
		p.NominalX(xtext[0], xtext[1], xtext[2], xtext[3], xtext[4])
	}

	if err := p.Save(6*vg.Inch, 15*vg.Inch, "./results/boxplot.pdf"); err != nil {
		panic(err)
	}

	f, err := os.Create("./results/boxplotStatistics.txt")
	check(err)
	defer f.Close()

	for _, str := range statitics {
		f.WriteString(str + "\n")
	}
	f.Sync()
	f.Close()
}