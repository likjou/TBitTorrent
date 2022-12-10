package tui

import (
	"fmt"
	"strconv"
	"time"

	"text/tabwriter"

	u "github.com/likjou/TBitTorrent/utils"

	human "github.com/dustin/go-humanize"
	"github.com/jroimartin/gocui"
	"github.com/olekukonko/tablewriter"
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

			//temp table specify view as io
			table2 := tablewriter.NewWriter(v)

			//temp table style
			table2.SetColWidth(55)
			table2.SetAutoWrapText(true)
			table2.SetBorder(false)
			table2.SetHeaderLine(false)
			table2.SetCenterSeparator("_")
			table2.SetColumnSeparator("|")
			table2.SetRowSeparator("_")

			//temp header
			table2.SetHeader([]string{"NAME", "SIZE", "PROGRESS", "STATUS", "SEEDS", "PEERS", "DOWN SPEED", "UP SPEED"})

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

				if len(torName) >= 55 {
					torName = torName[:50] + "..."
				}

				table2.Append([]string{newName, torSize, torProg, torStatus, torSeed, torPeers, torDownSpeed, torUpSpeed})
			}

			table2.Render()

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

		table2 := tablewriter.NewWriter(v2)

		//temp table style
		table2.SetAutoWrapText(false)
		table2.SetBorder(false)
		table2.SetHeaderLine(false)
		table2.SetCenterSeparator("_")
		table2.SetColumnSeparator("|")
		table2.SetRowSeparator("_")

		table2.Append([]string{"NAME", "SIZE", "PROGRESS", "STATUS", "SEEDS", "PEERS", "DOWN SPEED", "UP SPEED"})
		table2.Render()
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

		const padding = 0
		w := tabwriter.NewWriter(v4, 30, 0, padding, ' ', tabwriter.TabIndent)

		fmt.Fprintln(w, "Connection: \tETA: \tSeeds: \t")
		fmt.Fprintln(w, "Downloaded: \tUploaded: \tPeers: \t")
		fmt.Fprintln(w, "Down Speed: \tUp Speed: \tWasted: \t")

		w.Flush()
	}

	if v5, err := g.SetView("information", int(0.15*float32(maxX)+1), int(0.65*float32(maxY)+1), maxX-2, int(0.9*float32(maxY)-1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v5.Title = "Information"
		const padding = 0
		w := tabwriter.NewWriter(v5, 30, 0, padding, ' ', tabwriter.TabIndent)

		fmt.Fprintln(w, "Size: \tFiles: \tAdded at: \t")
		fmt.Fprintln(w, "InfoHash: \tPieces: \t")
		fmt.Fprintln(w, "Save Path: \t")

		w.Flush()

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

		table2 := tablewriter.NewWriter(v8)

		//temp table style
		table2.SetAutoWrapText(false)
		table2.SetBorder(false)
		table2.SetHeaderLine(false)
		table2.SetCenterSeparator("_")
		table2.SetColumnSeparator("|")
		table2.SetRowSeparator("_")

		//temp header
		table2.SetHeader([]string{"NO", "URL", "STATUS", "SEEDS", "LEECHERS"})
		table2.Render()

		v8.Title = "Tracker"
		g.SetViewOnBottom(v8.Name())
		v8.SetCursor(0, 1)

	}

	if v9, err := g.SetView("peersView", int(0.15*float32(maxX)+1), int(0.4*float32(maxY)+1), maxX-2, int(0.9*float32(maxY)-1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		table2 := tablewriter.NewWriter(v9)
		//temp table style
		table2.SetAutoWrapText(false)
		table2.SetBorder(false)
		table2.SetHeaderLine(false)
		table2.SetCenterSeparator("_")
		table2.SetColumnSeparator("|")
		table2.SetRowSeparator("_")

		//temp header
		table2.SetHeader([]string{"ADDRESS", "CLIENT", "DOWNLOAD SPEED", "UPLOAD SPEED", "SOURCE"})
		table2.Render()

		v9.Title = "Peers"
		g.SetViewOnBottom(v9.Name())
		v9.SetCursor(0, 1)
	}

	return nil
}
