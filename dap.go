package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

var (
	title       = flag.String("name", "", "set window title")
	exitOnWrite = flag.Bool("exit", false, "exit on files drppped")
)

const MaxDataSize = 256 * 64

func main() {
	flag.Parse()

	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle(*title)
	window.Connect("destroy", gtk.MainQuit)

	targets := []gtk.TargetEntry{
		{"text/uri-list", 0, 0},
		{"STRING", 0, 1},
		{"text/plain", 0, 2},
	}
	dest := gtk.NewLabel("drop me file")
	dest.DragDestSet(
		gtk.DEST_DEFAULT_MOTION|
			gtk.DEST_DEFAULT_HIGHLIGHT|
			gtk.DEST_DEFAULT_DROP,
		targets,
		gdk.ACTION_COPY)
	dest.DragDestAddUriTargets()

	closeButton := gtk.NewButtonWithLabel("Close")
	closeButton.Connect("clicked", func() {
		gtk.MainQuit()
	})

	statusBar := gtk.NewStatusbar()
	contextId := statusBar.GetContextId("go-gtk")
	statusBar.Push(contextId, "Drag and Drop files!")

	dest.Connect("drag-data-received", func(ctx *glib.CallbackContext) {
		sdata := gtk.NewSelectionDataFromNative(unsafe.Pointer(ctx.Args(3)))
		if sdata == nil {
			return
		}

		a := (*[MaxDataSize]byte)(sdata.GetData())
		filenames := strings.Split(string(a[0:sdata.GetLength()-1]), "\n")

		paths := []string{}
		for _, f := range filenames {
			path, _, err := glib.FilenameFromUri(f)
			if err != nil {
				continue
			}
			paths = append(paths, path)
		}

		for _, p := range paths {
			println(p)
		}

		statusBar.Push(contextId, fmt.Sprintf("%d path(s) printed", len(paths)))

		if len(paths) == 0 {
			return
		}

		if *exitOnWrite {
			gtk.MainQuit()
		}
	})

	vbox := gtk.NewVBox(false, 0)
	vbox.SetBorderWidth(5)
	vbox.PackStart(dest, true, true, 1)
	vbox.PackStart(closeButton, false, false, 0)
	vbox.PackStart(statusBar, false, false, 0)

	window.Add(vbox)

	window.SetSizeRequest(300, 200)
	window.ShowAll()
	gtk.Main()
}
