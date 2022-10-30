package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/rain/torrent"
	human "github.com/dustin/go-humanize"
	"github.com/gosuri/uitable"
	"github.com/jroimartin/gocui"
	c "github.com/likjou/TBitTorrent/config"
)

var (
	viewArr          = []string{"side", "torList", "generalBtn", "trackerBtn", "peersBtn"}
	active           = 0
	allTors          []*torrent.Torrent
	FilteredTors     []*torrent.Torrent
	torSes           *torrent.Session
	CurrInfo         *torrent.Torrent
	CurrTorListView  string
	CurrTorInfoView  string
	currTrackerLine  = 1
	currTorListLine  = 1
	currTorPeersLine = 1
)

//major functions are here including keybindings, logics and ui rendering

// to delete a given view name
func DelView(g *gocui.Gui, v *gocui.View, name string) error {
	g.DeleteView(name)
	prevActive := viewArr[active]
	g.SetCurrentView(prevActive)
	return nil
}

func DelViewCustom(g *gocui.Gui, v *gocui.View, name string, nxtvw string) error {
	g.DeleteView(name)
	g.SetCurrentView(nxtvw)
	return nil
}

// scrolling up logic for delview
func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()

		if cy == 1 {
			v.SetCursor(0, 1)
		} else {
			v.SetCursor(cx, cy-1)
		}

	}
	return nil
}

// scrolling down logic delview
func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()

		num := len(allTors)

		v.SetCursor(cx, cy+1)

		if cy == num {
			v.SetCursor(cx, cy)
			return nil
		}
	}
	return nil
}

// function to delete selected torrent in torrent client
func DelTorrent(g *gocui.Gui, v *gocui.View) error {
	allTorID := allTors

	_, cy := v.Cursor()

	nl := cy - 1
	if len(allTors) != 0 {
		if err := torSes.RemoveTorrent(allTorID[nl].ID()); err != nil {
			log.Panicln(err)
		}
		allTors = torSes.ListTorrents()
		FilteredTors = allTors
		DelView(g, v, "delTorrentView")
	}

	return nil
}

// to quit the application
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// to add torrent using torrent file to torrent client
func AddTorrent(g *gocui.Gui, v *gocui.View) error {
	filePath, err := getLine(g, v, v.Name())

	f, err1 := os.Open(filePath)
	t, err2 := torSes.AddTorrent(f, nil)

	if err != nil {
		return err
	}

	if err1 != nil {
		if err := ErrView(g, err1); err != nil {
			return err
		}
	}

	if err2 != nil {
		return err
	}

	allTors = append(allTors, t)
	defer f.Close()
	FilteredTors = allTors
	DelView(g, v, "addTorrentView")
	return nil
}

// to add torrent via magnet links to client, works with URI but wanted only magnet hehe
func AddMagnet(g *gocui.Gui, v *gocui.View) error {
	line := v.Buffer()
	link := strings.Trim(line, "\n\r")
	t, err2 := torSes.AddURI(link, nil)

	if err2 != nil {
		if err := ErrView2(g, err2); err != nil {
			return err
		}
	}

	if t != nil {
		allTors = append(allTors, t)
		FilteredTors = allTors
	} else {
		return nil
	}
	DelView(g, v, "addMagView")
	return nil
}

// cursor down logic for side bar, bruh
func CursorDownSide(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		v.SetCursor(cx, cy+1)

		if cy > 4 {
			v.SetCursor(0, cy)
		}
	}
	return nil
}

// cursor down logic for side bar, bruh
func CursorUpSide(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		v.SetCursor(cx, cy-1)

		if cy < 0 {
			v.SetCursor(0, 0)
		}
	}
	return nil
}

// tab to change focus to next view
func NextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := SetCurrentViewOnTop(g, name); err != nil {
		return err
	}

	active = nextIndex
	return nil
}

// set current filter view
func SetTorListView(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()

	switch cy {
	case 0:
		CurrTorListView = "alltorrent"
		newList := filterData()
		FilteredTors = newList
	case 1:
		CurrTorListView = "completed"
		newList := filterData()
		FilteredTors = newList
	case 2:
		CurrTorListView = "seeding"
		newList := filterData()
		FilteredTors = newList
	case 3:
		CurrTorListView = "downloading"
		newList := filterData()
		FilteredTors = newList
	case 4:
		CurrTorListView = "stopped"
		newList := filterData()
		FilteredTors = newList
	case 5:
		CurrTorListView = "verifying"
		newList := filterData()
		FilteredTors = newList
	}
	return nil
}

// create a new view to add torrent via torrent file
func AddTorrentView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	name := "addTorrentView"
	if v, err := g.SetView(name, int(0.1*float32(maxX)), int(0.45*float32(maxY)), int(0.9*float32(maxX)), int(0.55*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Clear()
		v.Title = "Torrent File Path: (Ctrl + a to cancel)"
		v.Editable = true
		v.Editor = gocui.DefaultEditor

		if _, err2 := g.SetCurrentView(name); err != nil {
			return err2
		}
	}
	return nil
}

// create a help manu to help user navigate
func Help(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	name := "help"
	v, err := g.SetView(name, int(0.25*float32(maxX)), int(0.25*float32(maxY)), int(0.75*float32(maxX)), int(0.75*float32(maxY)))
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Help Menu (Ctrl + h to close Help Menu)"
		fmt.Fprintf(v, "- Ctrl + a, to add torrent via torrent file\n")
		fmt.Fprintf(v, "- Ctrl + q, to add torrent via magnet link\n")
		fmt.Fprintf(v, "- Ctrl + d, to delete existing torrent\n")
		fmt.Fprintf(v, "- Ctrl + p, to stop/pause all running torrent\n")
		fmt.Fprintf(v, "- Ctrl + s, to start/resume all running torrent\n")
		fmt.Fprintf(v, "- Ctrl + h, to open help menu\n")
		fmt.Fprintf(v, "- Ctrl + c, to close application\n")
		fmt.Fprintf(v, "- Ctrl + t, to stop/pause a running torrent\n")
		fmt.Fprintf(v, "- Ctrl + y, to start a running torrent\n")
		fmt.Fprintf(v, "- Tab, to switch between views\n")

		if _, err := g.SetCurrentView(name); err != nil {
			return err
		}
	}
	return nil
}

// create a view to delete existing torrents
func DelTorrentView(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	name := "delTorrentView"
	if v, err := g.SetView(name, int(0.25*float32(maxX)), int(0.25*float32(maxY)), int(0.75*float32(maxX)), int(0.75*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Delete Torrent (Press Ctrl + d to cancel)"
		v.Highlight = true
		v.SetCursor(0, 1)
		fmt.Fprintf(v, "%v", "Existing Torrents:\n")

		for i, n := range allTors {
			fmt.Fprintf(v, "[%v] %v\n", i, n.Name())
		}
		if _, err2 := g.SetCurrentView(name); err != nil {
			return err2
		}
	}
	return nil
}

// to pause/stop all  torrents
func PauseTorrent(g *gocui.Gui, v *gocui.View) error {
	if err := torSes.StopAll(); err != nil {
		return err
	}
	return nil
}

// to start/resume all torrents
func StartTorrent(g *gocui.Gui, v *gocui.View) error {

	tors := allTors
	var allStatus []string

	for _, n := range tors {
		allStatus = append(allStatus, n.Stats().Status.String())
	}

	stat := contains(allStatus, "Stopping")

	if stat {
		if err := ErrViewMsg(g, "Please Wait For all torrent to Stop"); err != nil {
			return err
		}
	} else {
		if err := torSes.StartAll(); err != nil {
			return err
		}
	}
	return nil
}

// create a view to add torrent cia magnet links or URI
func AddMagnetView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("addMagView", int(0.1*float32(maxX)), int(0.45*float32(maxY)), int(0.9*float32(maxX)), int(0.65*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Clear()
		v.Title = "Torrent Magnet Link: (Ctrl + q to cancel)"
		v.Editable = true
		v.Editor = gocui.DefaultEditor
		v.Wrap = true

		if _, err2 := g.SetCurrentView("addMagView"); err != nil {
			return err2
		}
	}
	return nil
}

// scroll down logic for torrent list
func CursorDownTorList(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	_, vy := v.Size()

	if currTorListLine < len(allTors) {
		v.SetCursor(cx, cy+1)
		currTorListLine++
		if cy+1 == vy && currTorListLine <= len(allTors) {
			ox, oy := v.Origin()
			v.SetOrigin(ox, oy+1)
		}
	} else {
		v.SetCursor(cx, cy)
	}

	return nil
}

// scroll up logic for torrent list
func CursorUpTorList(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()

	if cy == 1 {
		if cy == 1 && currTorListLine != 1 {
			ox, oy := v.Origin()
			v.SetOrigin(ox, oy-1)
			currTorListLine--
		} else {
			v.SetCursor(cx, cy)
		}
	} else {
		v.SetCursor(cx, cy-1)
		currTorListLine--
	}
	return nil
}

// get torrent info for each selected torrent in torrent list
func GetTorInfo(g *gocui.Gui, v *gocui.View) error {
	if len(allTors) != 0 {
		CurrInfo = allTors[currTorListLine-1]
	} else {
		return nil
	}
	return nil
}

// pause/stop selected torrent
func PauseSingleTor(g *gocui.Gui, v *gocui.View) error {
	v2, _ := g.View("torList")

	_, cy := v2.Cursor()

	nl := cy - 1

	selTor := allTors[nl]

	err := selTor.Stop()
	if err != nil {
		log.Panicln(err)
	}
	return nil
}

// start/resume selected torrent
func StartSingleTor(g *gocui.Gui, v *gocui.View) error {
	v2, _ := g.View("torList")

	_, cy := v2.Cursor()

	nl := cy - 1

	selTor := allTors[nl]

	err := selTor.Start()
	if err != nil {
		log.Panicln(err)
	}
	return nil
}

// set the info view to general view
func SetGeneralView(g *gocui.Gui, v *gocui.View) error {
	CurrTorInfoView = "generalView"
	v1, err := g.View("transInfo")
	v2, err2 := g.View("information")
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}

	g.SetViewOnTop(v1.Name())
	g.SetViewOnTop(v2.Name())
	return nil
}

// set the info view to tracker view
func SetTrackerView(g *gocui.Gui, v *gocui.View) error {
	CurrTorInfoView = "trackerView"
	_, err := g.SetViewOnTop(CurrTorInfoView)
	_, err2 := g.SetCurrentView(CurrTorInfoView)
	if err != nil {
		log.Panicln(err)
	}
	if err2 != nil {
		log.Panicln(err2)
	}
	return nil
}

// set the info view to peers view
func SetPeerView(g *gocui.Gui, v *gocui.View) error {
	CurrTorInfoView = "peersView"
	_, err := g.SetViewOnTop(CurrTorInfoView)
	_, err2 := g.SetCurrentView(CurrTorInfoView)
	if err != nil {
		log.Panicln(err)
	}
	if err2 != nil {
		log.Panicln(err2)
	}
	return nil
}

// scroll down logic for tracker view
func CursorDownTracker(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	_, vy := v.Size()

	if currTrackerLine < len(CurrInfo.Trackers()) {
		v.SetCursor(cx, cy+1)
		currTrackerLine++
		if cy+1 == vy && currTrackerLine <= len(CurrInfo.Trackers()) {
			ox, oy := v.Origin()
			v.SetOrigin(ox, oy+1)
		}
	} else {
		v.SetCursor(cx, cy)
	}
	return nil
}

// scroll up logic for tracker view
func CursorUpTracker(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()

	if cy == 1 {
		if cy == 1 && currTrackerLine != 1 {
			ox, oy := v.Origin()
			v.SetOrigin(ox, oy-1)
			currTrackerLine--
		} else {
			v.SetCursor(cx, cy)
		}
	} else {
		v.SetCursor(cx, cy-1)
		currTrackerLine--
	}
	return nil
}

// scroll down logic gor peers view
func CursorDownPeers(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	_, vy := v.Size()

	if currTorListLine < len(CurrInfo.Peers()) {
		v.SetCursor(cx, cy+1)
		currTorListLine++
		if cy+1 == vy && currTorListLine <= len(CurrInfo.Peers()) {
			ox, oy := v.Origin()
			v.SetOrigin(ox, oy+1)
		}
	} else {
		v.SetCursor(cx, cy)
	}

	return nil
}

// scroll up logic for peers view
func CursorUpPeers(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()

	if cy == 1 {
		if cy == 1 && currTorPeersLine != 1 {
			ox, oy := v.Origin()
			v.SetOrigin(ox, oy-1)
			currTorPeersLine--
		} else {
			v.SetCursor(cx, cy)
		}
	} else {
		v.SetCursor(cx, cy-1)
		currTorPeersLine--
	}
	return nil
}

// get entered value from a given view
func getLine(g *gocui.Gui, v *gocui.View, name string) (string, error) {
	v1, err := g.View(name)
	if err != nil {
		log.Panicln(err)
		return "", err
	}
	line := v1.Buffer()
	dir := strings.Trim(line, "\n\r")
	return dir, nil
}

// error message view
func ErrView(g *gocui.Gui, err1 error) error {
	maxX, maxY := g.Size()
	if vErr, err := g.SetView("errorView", int(0.1*float32(maxX)), int(0.45*float32(maxY)), int(0.9*float32(maxX)), int(0.55*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vErr.Title = "ERROR"
		fmt.Fprintf(vErr, "%v\n", err1)
		fmt.Fprintln(vErr, "Press Enter to Continue.")

		if _, err2 := g.SetCurrentView(vErr.Name()); err != nil {
			return err2
		}
	}
	return nil
}

// error message view but no 2 LOL
func ErrView2(g *gocui.Gui, err1 error) error {
	maxX, maxY := g.Size()
	if vErr, err := g.SetView("errorView2", int(0.1*float32(maxX)), int(0.45*float32(maxY)), int(0.9*float32(maxX)), int(0.55*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vErr.Title = "ERROR"
		fmt.Fprintf(vErr, "%v\n", err1)
		fmt.Fprintln(vErr, "Press Enter to Continue.")

		if _, err2 := g.SetCurrentView(vErr.Name()); err != nil {
			return err2
		}
	}
	return nil
}

// error message but outputs string instead of error
func ErrViewMsg(g *gocui.Gui, msg string) error {
	maxX, maxY := g.Size()
	if vErr, err := g.SetView("errorViewMsg", int(0.1*float32(maxX)), int(0.45*float32(maxY)), int(0.9*float32(maxX)), int(0.55*float32(maxY))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vErr.Title = "ERROR"
		fmt.Fprintf(vErr, "%v\n", msg)
		fmt.Fprintln(vErr, "Press Enter to Continue.")

		if _, err2 := g.SetCurrentView(vErr.Name()); err != nil {
			return err2
		}
	}
	return nil
}

// set a given view on top
func SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

// filter torrent according to selected filter
func filterData() []*torrent.Torrent {

	var filteredList []*torrent.Torrent

	if strings.ToLower(CurrTorListView) == "completed" {
		for _, n := range allTors {
			if strings.ToLower(n.Stats().Status.String()) == "seeding" {
				filteredList = append(filteredList, n)
			}
		}
	} else if strings.ToLower(CurrTorListView) == "seeding" {
		for _, n := range allTors {
			if strings.ToLower(n.Stats().Status.String()) == "seeding" {
				filteredList = append(filteredList, n)
			}
		}
	} else if strings.ToLower(CurrTorListView) == "downloading" {
		for _, n := range allTors {
			if strings.ToLower(n.Stats().Status.String()) == "downloading" {
				filteredList = append(filteredList, n)
			}
		}
	} else if strings.ToLower(CurrTorListView) == "stopped" {
		for _, n := range allTors {
			if strings.ToLower(n.Stats().Status.String()) == "stopped" {
				filteredList = append(filteredList, n)
			}
		}
	} else if strings.ToLower(CurrTorListView) == "verifying" {
		for _, n := range allTors {
			if strings.ToLower(n.Stats().Status.String()) == "verifying" {
				filteredList = append(filteredList, n)
			}
		}
	} else {
		return allTors
	}
	return filteredList
}

// check if a string list contains a certain string
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// general info view for torrent info
func GeneralView(g *gocui.Gui) error {
	v4, err := g.View("transInfo")
	v5, err1 := g.View("information")

	if err != nil {
		panic(err)
	}
	if err1 != nil {
		panic(err1)
	}

	v4.Clear()
	v5.Clear()
	torInfo := CurrInfo

	transTable := uitable.New()
	transTable.MaxColWidth = 50
	transTable.RightAlign(0)
	transTable.RightAlign(1)
	transTable.RightAlign(2)
	transTable.RightAlign(3)
	transTable.RightAlign(4)
	transTable.RightAlign(5)

	infoTable := uitable.New()
	infoTable.MaxColWidth = 50
	infoTable.RightAlign(0)
	infoTable.RightAlign(1)
	infoTable.RightAlign(2)
	infoTable.RightAlign(3)
	infoTable.RightAlign(4)

	torPeers := strconv.Itoa(torInfo.Stats().Peers.Total)
	torInfoHash := torInfo.InfoHash().String()
	torDownSpeed := human.Bytes(uint64(torInfo.Stats().Speed.Download))
	torUpSpeed := human.Bytes(uint64(torInfo.Stats().Speed.Upload))
	torSize := human.Bytes(uint64(torInfo.Stats().Bytes.Total))

	torAddedAt := torInfo.AddedAt().Format(time.RFC822)
	torFileCount := strconv.Itoa(torInfo.Stats().FileCount)
	torWasted := human.Bytes(uint64(torInfo.Stats().Bytes.Wasted))
	torDownloaded := human.Bytes(uint64(torInfo.Stats().Bytes.Downloaded))
	torUploaded := human.Bytes(uint64(torInfo.Stats().Bytes.Uploaded))
	torConnection := strconv.Itoa(torInfo.Stats().Peers.Outgoing)
	pieces := strconv.Itoa(int(torInfo.Stats().Pieces.Total)) + " x " + human.Bytes(uint64(torInfo.Stats().PieceLength))
	savePath := c.AptorrentConfig.DataDir

	tracker := torInfo.Trackers()
	var totalSeed int
	for _, n := range tracker {
		totalSeed = totalSeed + n.Seeders
	}
	torSeed := strconv.Itoa(totalSeed)

	tempETA := torInfo.Stats().ETA
	var torETA string

	if tempETA == nil {
		torETA = "0"
	} else {
		torETA = tempETA.String()
	}
	transTable.AddRow("Connection: ", torConnection+"|", "ETA: ", torETA+"|", "Seeds: ", torSeed+"|")
	transTable.AddRow("Downloaded: ", torDownloaded+"|", "Uploaded: ", torUploaded+"|", "Peers: ", torPeers+"|")
	transTable.AddRow("Down Speed: ", torDownSpeed+"|", "Up Speed: ", torUpSpeed+"|", "Wasted: ", torWasted+"|")
	transTable.AddRow()
	fmt.Fprintf(v4, "%v", transTable)

	infoTable.AddRow("Size: ", torSize+"|", "Files: ", torFileCount+"|", "Added at: ", torAddedAt+"|")
	infoTable.AddRow("InfoHash: ", torInfoHash+"|", "Pieces: ", pieces+"|")
	infoTable.AddRow("Save Path: ", savePath+"|")
	fmt.Fprintf(v5, "%v", infoTable)
	return nil
}

// tracker view for torrent info
func TrackerView(g *gocui.Gui) error {
	statList := []string{"NotContactedYet", "Contacting", "Working", "Not Working"}
	v4, err := g.View("trackerView")
	if err != nil {
		panic(err)
	}

	v4.Clear()
	v4.Highlight = true

	tracker := CurrInfo.Trackers()
	table := uitable.New()
	table.MaxColWidth = 65
	table.AddRow("|NO", "|URL", "|STATUS", "|SEEDS", "|LEECHERS")

	for i, n := range tracker {

		table.AddRow("|"+strconv.Itoa(i+1), "|"+n.URL, "|"+statList[int(n.Status)], "|"+strconv.Itoa(n.Seeders), "|"+strconv.Itoa(n.Leechers))
	}
	fmt.Fprintf(v4, "%v", table)
	return nil
}

// peers view for torrent info
func PeersView(g *gocui.Gui) error {

	v4, err := g.View("peersView")
	if err != nil {
		panic(err)
	}
	v4.Clear()
	v4.Highlight = true

	allpeers := CurrInfo.Peers()
	table := uitable.New()
	table.MaxColWidth = 65
	table.AddRow("|ADDRESS", "|CLIENT", "|DOWNLOAD SPEED", "|UPLOAD SPEED", "|SOURCE")

	for _, n := range allpeers {
		table.AddRow("|"+n.Addr.String(), "|"+n.Client, "|"+human.Bytes(uint64(n.DownloadSpeed)), "|"+human.Bytes(uint64(n.UploadSpeed)), "|"+strconv.Itoa(int(n.Source)))
	}
	fmt.Fprintf(v4, "%v", table)
	return nil
}

// initialize the torrent session
func InitTorSess() error {
	// Create a session
	ses, err := torrent.NewSession(torrent.Config(c.AptorrentConfig))
	allTors = ses.ListTorrents()
	FilteredTors = ses.ListTorrents()
	torSes = ses
	if err != nil {
		return err
	}
	return nil
}
