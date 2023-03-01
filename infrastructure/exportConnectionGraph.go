package infrastructure

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// SaveConnectionGraphToFsAdapter save the actual Connection graph as png to FS
type SaveConnectionGraphToFsAdapter struct {
	FileAndFolder string
}

type SaveNodeGraphToFsAdapter struct {
	FileAndFolder string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// ExportConnectionGraph export execution
func (e SaveConnectionGraphToFsAdapter) ExportConnectionGraph(conncetionGraph string) {
	prefix_filename := e.FileAndFolder
	f, err := os.Create(prefix_filename + "networkgraph.gv")
	check(err)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	w := bufio.NewWriter(f)
	_, err = w.WriteString(conncetionGraph)

	if err != nil {
		return
	}
	err = w.Flush()
	if err != nil {
		fmt.Println("Error Export Connection Graph")
		return
	}

	path, _ := exec.LookPath("dot")
	cmd, _ := exec.Command(path, "-Tpng", prefix_filename+"networkgraph.gv").Output()
	mode := 0777
	err = os.WriteFile(prefix_filename+"networkgraph.png", cmd, os.FileMode(mode))
	if err != nil {
		return
	}
}

func (e SaveNodeGraphToFsAdapter) ExportNodeGraph(nodes []string) {
	prefix_filename := e.FileAndFolder
	f, err := os.Create(prefix_filename + "nodes.plantuml")
	check(err)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	w := bufio.NewWriter(f)

	for _, value := range nodes {
		_, err = w.WriteString(value)
		if err != nil {
			return
		}
	}
	err = w.Flush()
	if err != nil {
		fmt.Println("Error Export Node Graph")
		return
	}
}
