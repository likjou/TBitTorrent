package main

import (
	"log"

	k "github.com/likjou/TBitTorrent/keybindings"

	t "github.com/likjou/TBitTorrent/tui"

	u "github.com/likjou/TBitTorrent/utils"

	"github.com/cenkalti/rain/torrent"
	"github.com/jroimartin/gocui"
)

// run application
func runApp() {
	g, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(t.Layout)
	g.Highlight = true
	g.Cursor = false
	g.SelFgColor = gocui.ColorGreen

	if err := k.Keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

//main function

func main() {

	torrent.DisableLogging()
	err1 := u.InitTorSess()

	if err1 != nil {
		log.Panicln(err1)
	} else {
		runApp()
	}

}
