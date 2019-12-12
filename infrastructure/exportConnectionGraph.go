package infrastructure

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// SaveConnectionGraphToFsAdapter save the actual Connection graph as png to FS
type SaveConnectionGraphToFsAdapter struct {
	FileAndFolder string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// ExportConnectionGraph export excecution
func (e SaveConnectionGraphToFsAdapter) ExportConnectionGraph(conncetionGraph string) {
	filename := e.FileAndFolder
	f, err := os.Create(filename)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(conncetionGraph)
	w.Flush()

	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-Tpng", filename).Output()
	mode := int(0777)
	ioutil.WriteFile("./results/networkgraph.pdf", cmd, os.FileMode(mode))
	//dot -Tpng  > test.png && open test.png
	fmt.Println("        Networkgraph created")
}
