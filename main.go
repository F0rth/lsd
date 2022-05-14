package main

import (
	"bufio"
	"log"
	"os/exec"
	"strings"

	g "github.com/AllenDang/giu"
)

var (
	url, member_question, peasant_question string
	splitLayoutPos                         float32 = 320
	size                                   int32   = 10
)

func reset() {
	member_question = ""
	peasant_question = ""
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
			if strings.Contains(live, "Member") || strings.Contains(live, "New member") {
				member_question += live + "\n"
			} else {
				peasant_question += live + "\n"
			}
		}
	}
}
func loop() {

	g.SingleWindow().Layout(
		g.Row(
			g.Label("url:"),
			g.InputText(&url).Size(g.Auto-150),
			g.Button("scrap!").OnClick(scrap),
			g.Button("reset").OnClick(reset),
		),
		g.SliderInt(&size, 6, 40).Size(g.Auto),
		g.SplitLayout(g.DirectionVertical, &splitLayoutPos,
			// g.CodeEditor().Text(member_question).Border(true).ShowWhitespaces(false).TabSize(2),
			// g.CodeEditor().Text(peasant_question).Border(true).ShowWhitespaces(false).TabSize(2),
			g.Style().SetFontSize(float32(size)).To(g.InputTextMultiline(&member_question).Size(g.Auto, g.Auto)),
			g.Style().SetFontSize(float32(size)).To(g.InputTextMultiline(&peasant_question).Size(g.Auto, g.Auto)),
		),
	)
}

func main() {
	wnd := g.NewMasterWindow("Live Shill Discriminator", 800, 800, g.MasterWindowFlagsFloating)
	wnd.Run(loop)

}
