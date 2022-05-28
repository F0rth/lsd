package main

import (
	"bufio"
	//"fmt"

	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

var (
	url                       string
	binary                    string  = "chat_downloader"
	options                   string  = "messages, superchat, tickers, banners, donations, engagement, purchases"
	splitLayoutPos            float32 = 320
	size                      int32   = 10
	pasta                     bool    = false
	editor1, editor2, editor3 *g.CodeEditorWidget
	f1, f2, f3                *os.File
	err1, err2, err3          error
	errMarkers                imgui.ErrorMarkers
	cmd                       *exec.Cmd
)

func saveTofile() {
	os.WriteFile(f1.Name(), []byte(editor1.GetText()), 0666)
	os.WriteFile(f2.Name(), []byte(editor2.GetText()), 0666)
	os.WriteFile(f3.Name(), []byte(editor2.GetText()), 0666)

}

func reset() {
	editor1.Text("")
	editor2.Text("")
	editor3.Text("")
}

func scrap() {
	go chat_download()
}

func rm() {
	os.Remove(f1.Name())
	os.Remove(f2.Name())
	os.Remove(f3.Name())
}

func kill() {
	cmd.Process.Kill()
}

func chat_download() {
	cmd = exec.Command(binary, "--message_groups", options, url)
	//fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println("start")
	err = cmd.Start()
	//fmt.Println(cmd.Process.Pid)

	if err != nil {
		log.Panic(err)
	}
	buf := bufio.NewReader(stdout)
	for {
		charset := "Windows-1252"
		line, _, _ := buf.ReadLine()
		live := string(line)
		result2 := live
		// windows charset
		if runtime.GOOS == "windows" {
			e, err := ianaindex.MIME.Encoding(charset)
			if err != nil {
				log.Fatal(err)
			}
			r := transform.NewReader(bytes.NewBufferString(live), e.NewDecoder())
			result, err := ioutil.ReadAll(r)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Printf("%s\n", result)
			result2 = string(result)
		}

		if strings.Contains(result2, "*€") || strings.Contains(result2, "*$") || strings.Contains(result2, "*£") {
			editor3.InsertText(result2 + "\n")
			g.Update()
			//errMarkers.Clear()
			//errMarkers.Insert(1, "Error message")
			//editor1.ErrorMarkers(errMarkers)
		} else if strings.Contains(result2, "?") {
			if strings.Contains(result2, "Member") || strings.Contains(result2, "New member") || strings.Contains(result2, "Moderator") {
				editor1.InsertText(result2 + "\n")
				g.Update()
			} else {
				editor2.InsertText(result2 + "\n")
				g.Update()
			}
		}
		saveTofile()
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
			g.Button("rm").OnClick(rm),
			g.Checkbox("pasta", &pasta),
		),
		g.Row(
			g.Label("binary path:"),
			g.InputText(&binary),
			g.Button("kill").OnClick(kill),
		),
		g.Row(
			g.Label("options"),
			g.InputText(&options),
		),
		g.SliderInt(&size, 6, 40).Size(g.Auto),
		g.SplitLayout(split_direction, &splitLayoutPos,
			g.Layout{
				g.Style().SetFontSize(float32(size)).To(editor3),
			},
			g.SplitLayout(split_direction, &splitLayoutPos,
				g.Style().SetFontSize(float32(size)).To(editor1),
				g.Style().SetFontSize(float32(size)).To(editor2),
			),
		),
	)
}

func main() {
	if runtime.GOOS == "windows" {
		binary = ".\\chat_downloader.exe"
	}
	f1name := "members.txt"
	f2name := "nonmembers.txt"
	f3name := "superchat.txt"
	editor1 = g.CodeEditor().Border(true).ShowWhitespaces(false).TabSize(2)
	editor2 = g.CodeEditor().Border(true).ShowWhitespaces(false).TabSize(2)
	editor3 = g.CodeEditor().Border(true).ShowWhitespaces(false).TabSize(2)
	f1, err1 = os.Open(f1name)
	if os.IsNotExist(err1) {
		f1, _ = os.Create(f1name)
	}
	file1, _ := os.ReadFile(f1.Name())
	editor1.InsertText(string(file1))

	f2, err2 = os.Open(f2name)
	if os.IsNotExist(err2) {
		f2, _ = os.Create(f2name)
	}
	file2, _ := os.ReadFile(f2.Name())
	editor2.InsertText(string(file2))

	f3, err3 = os.Open(f3name)
	if os.IsNotExist(err3) {
		f2, _ = os.Create(f3name)
	}
	file3, _ := os.ReadFile(f3.Name())
	editor3.InsertText(string(file3))

	defer f1.Close()
	defer f2.Close()
	defer f3.Close()

	wnd := g.NewMasterWindow("Live Shill Discriminator", 800, 800, g.MasterWindowFlagsFloating)
	wnd.Run(loop)

}
