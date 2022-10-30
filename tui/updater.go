package tui

import (
	"fmt"
	"strconv"
	"time"

	human "github.com/dustin/go-humanize"
	"github.com/gosuri/uitable"
	"github.com/jroimartin/gocui"
	u "github.com/likjou/TBitTorrent/utils"
)

const (
	MIN_WIDTH  = 113
	MIN_HEIGHT = 31
)

//layouts and updater functions are here

// torrent list updater with filter
func torListTicker(g *gocui.Gui, v *gocui.View) {

	ticker := time.NewTicker(1 * time.Second)

	for _ = range ticker.C {

		v.Clear()
		g.Update(func(g *gocui.Gui) error {
			v.Title = "(Press Ctrl + h for help)"
			v.Highlight = true

			table := uitable.New()
			table.MaxColWidth = 50
			table.AddRow("|NAME", "|SIZE", "|PROGRESS", "|STATUS", "|SEEDS", "|PEERS", "|DOWN SPEED", "|UP SPEED")

			for _, tor := range u.FilteredTors {
				torName := tor.Name()
				torSize := human.Bytes(uint64(tor.Stats().Bytes.Total))
				torProg := human.Bytes(uint64(tor.Stats().Bytes.Completed))
				torStatus := tor.Stats().Status.String()
				trackers := tor.Trackers()
				totalSeed := 0
				for _, n := range trackers {
					totalSeed = totalSeed + n.Seeders
				}
				torSeed := strconv.Itoa(totalSeed)
				torPeers := strconv.Itoa(tor.Stats().Peers.Total)
				torDownSpeed := human.Bytes(uint64(tor.Stats().Speed.Download))
				torUpSpeed := human.Bytes(uint64(tor.Stats().Speed.Upload))

				table.AddRow("|"+torName, "|"+torSize, "|"+torProg, "|"+torStatus, "|"+torSeed, "|"+torPeers, "|"+torDownSpeed, "|"+torUpSpeed)
			}

			fmt.Fprintf(v, "%v", table)

			return nil
		})
	}
}

// torrent info updater for each selected torrent
func transInfoTicker(g *gocui.Gui) {

	ticker := time.NewTicker(1 * time.Second)

	for _ = range ticker.C {
		g.Update(func(g *gocui.Gui) error {

			switch u.CurrTorInfoView {
			case "generalView":
				u.GeneralView(g)
			case "trackerView":
				u.TrackerView(g)
			case "peersView":
				u.PeersView(g)
			}

			return nil
		})
	}

	if u.CurrTorInfoView == "peersView" {
		v4, err := g.View("peersView")
		if err != nil {
			panic(err)
		}
		g.SetViewOnTop(v4.Name())
	} else if u.CurrTorInfoView == "trackerView" {
		v4, err := g.View("trackerView")
		if err != nil {
			panic(err)
		}

		g.SetViewOnTop(v4.Name())
	} else {
		v4, err := g.View("generalView")
		if err != nil {
			panic(err)
		}

		g.SetViewOnTop(v4.Name())
	}
}

// function to check if a torrent is selected from torListView
func checkCurrInfo(g *gocui.Gui) {
	ticker := time.NewTicker(1 * time.Second)

	for _ = range ticker.C {

		if u.CurrInfo != nil {
			ticker.Stop()
			go transInfoTicker(g)
		}
	}
}

// Layout of the entire terminal ui
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if maxX < MIN_WIDTH || maxY < MIN_HEIGHT {
		if v, err := g.SetView("errDim", 0, 0, maxX, maxY); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Frame = true
			v.Title = "Error"
			g.Cursor = false
			fmt.Fprintln(v, "Terminal is too small")
		}
		return nil
	}
	if _, err := g.View("errDim"); err == nil {
		g.DeleteView("errDim")
	}

	if v1, err := g.SetView("side", 0, 0, int(0.15*float32(maxX)-1), maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v1, "All Torrents")
		fmt.Fprintln(v1, "Completed")
		fmt.Fprintln(v1, "Seeding")
		fmt.Fprintln(v1, "Downloading")
		fmt.Fprintln(v1, "Stopped")
		fmt.Fprintln(v1, "Verifying")

		v1.Title = "TBitTorrent"
		v1.Highlight = true

		if _, err = u.SetCurrentViewOnTop(g, "side"); err != nil {
			return err
		}

	}
	if v2, err := g.SetView("torList", int(0.15*float32(maxX)), 0, maxX-1, int(0.4*float32(maxY)-1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v2.Highlight = true
		v2.SetCursor(0, 1)
		v2.Title = "(Press Ctrl + h for help)"

		table := uitable.New()
		table.MaxColWidth = 50
		table.AddRow("|NAME", "|SIZE", "|PROGRESS", "|STATUS", "|SEEDS", "|PEERS", "|DOWN SPEED", "|UP SPEED")
		fmt.Fprintf(v2, "%v", table)

		u.CurrTorListView = "alltorrent"

		go torListTicker(g, v2)

	}
	if _, err := g.SetView("torInfo", int(0.15*float32(maxX)), int(0.4*float32(maxY)), maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		u.CurrTorInfoView = "generalView"
		go checkCurrInfo(g)
	}

	if v4, err := g.SetView("transInfo", int(0.15*float32(maxX)+1), int(0.4*float32(maxY)+1), maxX-2, int(0.65*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v4.Title = "Transfer"

		transTable := uitable.New()
		transTable.MaxColWidth = 50
		transTable.RightAlign(0)
		transTable.RightAlign(2)
		transTable.RightAlign(4)

		transTable.AddRow("Connection: ", " |", "ETA: ", " |", "Seeds: ", " |")
		transTable.AddRow("Downloaded: ", " |", "Uploaded: ", " |", "Peers: ", " |")
		transTable.AddRow("Down Speed: ", " |", "Up Speed: ", " |", "Wasted: ", " |")
		fmt.Fprintf(v4, "%v", transTable)
	}

	if v5, err := g.SetView("information", int(0.15*float32(maxX)+1), int(0.65*float32(maxY)+1), maxX-2, int(0.9*float32(maxY)-1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v5.Title = "Information"
		infoTable := uitable.New()
		infoTable.MaxColWidth = 50
		infoTable.RightAlign(0)
		infoTable.RightAlign(2)

		infoTable.AddRow("Size: ", " |", "Files: ", " |", "Added at:", " |")
		infoTable.AddRow("InfoHash: ", " |", "Pieces: ", " |")
		infoTable.AddRow("Save Path: ", " |")
		fmt.Fprintf(v5, "%v", infoTable)

	}

	if v6, err := g.SetView("generalBtn", int(0.15*float32(maxX)+1), int(0.9*float32(maxY)), int(0.25*float32(maxX)), int(0.95*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v6, "%v", "General")

	}

	if v7, err := g.SetView("trackerBtn", int(0.25*float32(maxX)+1), int(0.9*float32(maxY)), int(0.35*float32(maxX)), int(0.95*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v7, "%v", "Tracker")

	}

	if v7, err := g.SetView("peersBtn", int(0.35*float32(maxX)+1), int(0.9*float32(maxY)), int(0.45*float32(maxX)), int(0.95*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v7, "%v", "Peers")

	}

	if v8, err := g.SetView("trackerView", int(0.15*float32(maxX)+1), int(0.4*float32(maxY)+1), maxX-2, int(0.9*float32(maxY)-1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		table := uitable.New()
		table.MaxColWidth = 65
		table.AddRow("|NO", "|URL", "|STATUS", "|SEEDS", "|LEECHERS")
		fmt.Fprintf(v8, "%v", table)
		v8.Title = "Tracker"
		g.SetViewOnBottom(v8.Name())
		v8.SetCursor(0, 1)

	}

	if v9, err := g.SetView("peersView", int(0.15*float32(maxX)+1), int(0.4*float32(maxY)+1), maxX-2, int(0.9*float32(maxY)-1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		table := uitable.New()
		table.MaxColWidth = 65
		table.AddRow("|ADDRESS", "|CLIENT", "|DOWNLOAD SPEED", "|UPLOAD SPEED", "|SOURCE")
		fmt.Fprintf(v9, "%v", table)
		v9.Title = "Peers"
		g.SetViewOnBottom(v9.Name())
		v9.SetCursor(0, 1)
	}

	return nil
}
