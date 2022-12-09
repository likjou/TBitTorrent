package keybindings

import (
	"log"

	u "aptorrent/utils"

	"github.com/jroimartin/gocui"
)

// all key presses for each individual view
// really ugly but oh well ;)

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, u.Quit); err != nil {
		log.Panicln(err)
	}

	for _, n := range []string{"delTorrentView"} {
		name := "delTorrentView"
		if err := g.SetKeybinding(n, gocui.KeyCtrlD, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelView(g, v, name)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyArrowUp, gocui.ModNone, u.CursorUp); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyArrowDown, gocui.ModNone, u.CursorDown); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyEnter, gocui.ModNone, u.DelTorrent); err != nil {
			return err
		}

	}

	for _, n := range []string{"errorView"} {

		if err := g.SetKeybinding(n, gocui.KeyEnter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelViewCustom(g, v, v.Name(), "addTorrentView")
			}); err != nil {
			return err
		}
	}

	for _, n := range []string{"errorView2"} {

		if err := g.SetKeybinding(n, gocui.KeyEnter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelViewCustom(g, v, v.Name(), "addMagView")
			}); err != nil {
			return err
		}
	}

	for _, n := range []string{"errorViewMsg"} {
		if err := g.SetKeybinding(n, gocui.KeyEnter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelView(g, v, v.Name())
			}); err != nil {
			return err
		}
	}

	for _, n := range []string{"addTorrentView"} {
		if err := g.SetKeybinding(n, gocui.KeyCtrlA, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelView(g, v, v.Name())
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyEnter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddTorrent(g, v)
			}); err != nil {
			return err
		}
	}

	for _, n := range []string{"addMagView"} {
		if err := g.SetKeybinding(n, gocui.KeyCtrlQ, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelView(g, v, v.Name())
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyEnter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddMagnet(g, v)
			}); err != nil {
			return err
		}

	}

	for _, n := range []string{"help"} {
		name := "help"
		if err := g.SetKeybinding(n, gocui.KeyCtrlH, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelView(g, v, name)
			}); err != nil {
			return err
		}

	}

	for _, n := range []string{"side"} {
		if err := g.SetKeybinding(n, gocui.KeyArrowUp, gocui.ModNone, u.CursorUpSide); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyArrowDown, gocui.ModNone, u.CursorDownSide); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyTab, gocui.ModNone, u.NextView); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(n, gocui.KeyEnter, gocui.ModNone, u.SetTorListView); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlA, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddTorrentView(g)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlH, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.Help(g)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlD, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelTorrentView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlP, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.PauseTorrent(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlS, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.StartTorrent(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlQ, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddMagnetView(g)
			}); err != nil {
			return err
		}

	}

	for _, v := range []string{"torList"} {

		if err := g.SetKeybinding(v, gocui.KeyArrowUp, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.CursorUpTorList(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(v, gocui.KeyArrowDown, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.CursorDownTorList(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(v, gocui.KeyEnter, gocui.ModNone, u.GetTorInfo); err != nil {
			return err
		}

		if err := g.SetKeybinding(v, gocui.KeyTab, gocui.ModNone, u.NextView); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(v, gocui.KeyCtrlT, gocui.ModNone, u.PauseSingleTor); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(v, gocui.KeyCtrlY, gocui.ModNone, u.StartSingleTor); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(v, gocui.KeyCtrlA, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddTorrentView(g)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(v, gocui.KeyCtrlH, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.Help(g)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(v, gocui.KeyCtrlD, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelTorrentView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(v, gocui.KeyCtrlP, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.PauseTorrent(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(v, gocui.KeyCtrlS, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.StartTorrent(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(v, gocui.KeyCtrlQ, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddMagnetView(g)
			}); err != nil {
			return err
		}
	}

	for _, n := range []string{"generalBtn"} {
		if err := g.SetKeybinding(n, gocui.KeyTab, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.NextView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyEnter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.SetGeneralView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlA, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddTorrentView(g)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlH, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.Help(g)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlD, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelTorrentView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlP, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.PauseTorrent(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlS, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.StartTorrent(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlQ, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddMagnetView(g)
			}); err != nil {
			return err
		}

	}

	for _, n := range []string{"trackerBtn"} {
		if err := g.SetKeybinding(n, gocui.KeyTab, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.NextView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyEnter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.SetTrackerView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlA, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddTorrentView(g)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlH, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.Help(g)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlD, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelTorrentView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlP, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.PauseTorrent(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlS, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.StartTorrent(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlQ, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddMagnetView(g)
			}); err != nil {
			return err
		}
	}

	for _, n := range []string{"peersBtn"} {
		if err := g.SetKeybinding(n, gocui.KeyTab, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.NextView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyEnter, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.SetPeerView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlA, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddTorrentView(g)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlH, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.Help(g)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlD, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.DelTorrentView(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlP, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.PauseTorrent(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlS, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.StartTorrent(g, v)
			}); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyCtrlQ, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.AddMagnetView(g)
			}); err != nil {
			return err
		}
	}

	for _, n := range []string{"trackerView"} {
		if err := g.SetKeybinding(n, gocui.KeyArrowUp, gocui.ModNone, u.CursorUpTracker); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyArrowDown, gocui.ModNone, u.CursorDownTracker); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyTab, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.NextView(g, v)
			}); err != nil {
			return err
		}

	}

	for _, n := range []string{"peersView"} {
		if err := g.SetKeybinding(n, gocui.KeyArrowUp, gocui.ModNone, u.CursorUpPeers); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyArrowDown, gocui.ModNone, u.CursorDownPeers); err != nil {
			return err
		}

		if err := g.SetKeybinding(n, gocui.KeyTab, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				return u.NextView(g, v)
			}); err != nil {
			return err
		}

	}
	return nil
}
