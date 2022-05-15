package main

import (
	"bufio"
	"log"
	"os/exec"
	"strings"

	g "github.com/AllenDang/giu"
)

var (
	url              string
	splitLayoutPos   float32 = 320
	size             int32   = 10
	pasta            bool    = false
	editor1, editor2 *g.CodeEditorWidget
)

func reset() {
	editor1.Text("")
	editor2.Text("")
}

func scrap() {
	go chat_download()
}

func chat_download() {
	cmd := exec.Command("chat_downloader", url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	buf := bufio.NewReader(stdout)
	for {
		line, _, _ := buf.ReadLine()
		live := string(line)
		if strings.Contains(live, "?") {
			if strings.Contains(live, "Member") || strings.Contains(live, "New member") || strings.Contains(live, "Moderator") {
				editor1.InsertText(live + "\n")
			} else {
				editor2.InsertText(live + "\n")
			}
		}
	}
}
func loop() {
	var split_direction g.SplitDirection
	if pasta {
		split_direction = g.DirectionHorizontal
	} else {
		split_direction = g.DirectionVertical
	}
	g.SingleWindow().Layout(
		g.Row(
			g.Label("url:"),
			g.InputText(&url).Size(g.Auto-300),
			g.Button("scrap!").OnClick(scrap),
			g.Button("oh crap").OnClick(reset),
			g.Checkbox("pasta", &pasta),
		),
		g.SliderInt(&size, 6, 40).Size(g.Auto),
		g.SplitLayout(split_direction, &splitLayoutPos,
			g.Style().SetFontSize(float32(size)).To(editor1),
			g.Style().SetFontSize(float32(size)).To(editor2),
		),
	)
}

func main() {
	editor1 = g.CodeEditor().Border(true).ShowWhitespaces(false).TabSize(2)
	editor2 = g.CodeEditor().Border(true).ShowWhitespaces(false).TabSize(2)
	wnd := g.NewMasterWindow("Live Shill Discriminator", 800, 800, g.MasterWindowFlagsFloating)
	wnd.Run(loop)

}
