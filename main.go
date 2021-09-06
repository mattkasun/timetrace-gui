package main

import (
	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/fs"
	"github.com/gin-gonic/gin"
)

//func status(table *tview.Table, timetrace *core.Timetrace) {
//	i := 0
//	for {
//		log.Print(i, string(i))
//		report, err := timetrace.Status()
//		if err != nil {
//			log.Fatal(err)
//		}
//		cell := tview.NewTableCell(report.Current.Project.Key)
//		table.SetCell(1, 0, cell)
//		table.SetCell(1, 1, tview.NewTableCell(timetrace.Formatter().FormatCurrentTime(report)))
//		table.SetCell(1, 2, tview.NewTableCell(timetrace.Formatter().FormatTodayTime(report)))
//		table.SetCell(1, 3, tview.NewTableCell(timetrace.Formatter().FormatBreakTime(report)))
//		table.SetCellSimple(2, 3, string(i))
//		if 1%2 == 0 {
//			table.SetBorders(true)
//		} else {
//			table.SetBorders(false)
//		}
//		time.Sleep(1 * time.Second)
//		i++
//	}
//}

var timetrace *core.Timetrace

func main() {
	config := config.Get()
	file := fs.New(config)
	timetrace = core.New(config, file)
	router := SetupRouter()
	router.Run("127.0.0.1:8080")
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("html/*")
	router.GET("/", DisplayLanding)
	router.POST("/", StartStop)

	return router
}

//	header := []string{"Project", "Since Start", "Today", "Breaks"}
//	config := config.Get()
//	fs := fs.New(config)
//
//	timetrace := core.New(config, fs)
//
//	app := tview.NewApplication()
//	app.EnableMouse(true)
//
//	//Set up Status Table
//	table := tview.NewTable()
//	table.SetBorders(true)
//	for i, name := range header {
//		cell := tview.NewTableCell(name)
//		cell.SetBackgroundColor(tcell.ColorBlue.TrueColor())
//		table.SetCell(0, i, cell)
//	}
//	report, err := timetrace.Status()
//	if err != nil {
//		log.Fatal(err)
//	}
//	cell := tview.NewTableCell(report.Current.Project.Key)
//	table.SetCell(1, 0, cell)
//	table.SetCell(1, 1, tview.NewTableCell(timetrace.Formatter().FormatCurrentTime(report)))
//	table.SetCell(1, 2, tview.NewTableCell(timetrace.Formatter().FormatTodayTime(report)))
//	table.SetCell(1, 3, tview.NewTableCell(timetrace.Formatter().FormatBreakTime(report)))
//
//	app.SetRoot(table, true)
//	//update status table
//	go status(table, timetrace)
//	err = app.Run()
//	if err != nil {
//		log.Fatal(err)
//	}
//}
