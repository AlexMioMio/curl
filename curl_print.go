package curl

import (
	"fmt"
	"strings"
	"time"
)

/*
 curl.Print Options
*/
type PrintOps struct {
	Header   bool
	Footer   bool
	LeftEnd  string
	RightEnd string
	Fill     string
	Arrow    string
	Empty    string
}

/*
 Set PrintOps default values
*/
var Options = PrintOps{true, true, "[", "]", "=", ">", "-"}

/*
 Print Header
*/
func header(dl *Download) {
	if Options.Header {
		fmt.Printf("Start download [%v].\n", strings.Join((*dl).GetValues("Title"), ", "))
	}
}

/*
 Print Footer
*/

func footer() {
	if Options.Footer {
		fmt.Println("\r--------\nEnd download.")
	}
}

/*
 title: 70% [==============>__________________] 925ms
*/
func progressbar(title string, start time.Time, i int, suffix string) {
	h := Options.LeftEnd + strings.Repeat(Options.Fill, i) + Options.Arrow + strings.Repeat(Options.Empty, 50-i) + Options.RightEnd
	d := time.Now().Sub(start)
	s := fmt.Sprintf("%v %.0f%% %s %v", safeTitle(title), float32(i)/50*100, h, time.Duration(d.Seconds())*time.Second)
	if len(s) > 80 {
		s = s[:80]
	}
	e := strings.Repeat(" ", 80-len(s))
	fmt.Printf("\r%v%v%v", s, e, suffix)
}