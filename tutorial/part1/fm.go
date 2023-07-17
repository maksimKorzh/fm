package main

import (
  "os"
  "fmt"
  "github.com/nsf/termbox-go"
)

var DRIVES = ""
var drive = 0
var COLS, ROWS int
var LEFT = 0
var RIGHT = 1
var active_panel = LEFT
var split_ch = "/"

type Panel struct {
  path string
  panel []map[string]interface{}
  current_row, row_offset int
}

var panels = []Panel {
  {
    path: "",
    panel: []map[string]interface{}{},
    current_row: 0,
    row_offset: 0,
  },
  {
    path: "",
    panel: []map[string]interface{}{},
    current_row: 0,
    row_offset: 0,
  },
}

func display_cell(lt, rt, lb, rb [] int) {
  for row := lt[1]+1; row < lb[1]; row++ {
    termbox.SetChar(lt[0], row, 9474)
    termbox.SetChar(rt[0], row, 9474)
    for col := lt[0]+1; col < rt[0]; col++ {
      termbox.SetChar(col, rt[1], 9472)
      termbox.SetChar(col, rb[1], 9472)
  }};termbox.SetChar(lt[0], lt[1], 9484)
  termbox.SetChar(lb[0], lb[1], 9492)
  termbox.SetChar(rt[0], rt[1], 9488)
  termbox.SetChar(rb[0], rb[1], 9496)
}

func display_borders() {
  llt := []int{0,0}
  lrt := []int{int(COLS/2)-1,0}
  llb := []int{0,ROWS-1}
  lrb := []int{int(COLS/2)-1,ROWS-1}
  rlt := []int{llt[0]+lrt[0]+1, llt[1]}
  rrt := []int{lrt[0]+lrt[0]+1, lrt[1]}
  rlb := []int{llb[0]+lrt[0]+1, llb[1]}
  rrb := []int{lrb[0]+lrt[0]+1, lrb[1]}
  display_cell(llt, lrt, llb, lrb)
  display_cell(rlt, rrt, rlb, rrb)
}

func get_event() termbox.Event {
  var event termbox.Event
  switch poll_event := termbox.PollEvent(); event.Type {
    case termbox.EventKey: event = poll_event
    case termbox.EventError: panic(event.Err)
  }; return event
}

func process_keypress() {
  event := get_event()
  switch event.Key {
    case termbox.KeyEsc: {
      termbox.Close();
      os.Exit(0)
    }
  }
}

func run_file_manager() {
  err := termbox.Init()
  if err != nil { fmt.Println(err); os.Exit(1) }
  panels[LEFT].path, err = os.Getwd()
  panels[RIGHT].path, err = os.Getwd()
  if err != nil { termbox.Close(); os.Exit(1) }
  for {
    COLS, ROWS = termbox.Size()
    if COLS < 78 { COLS = 78 }
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
    display_borders()
    termbox.Flush()
    process_keypress()
  }
}

func main() { run_file_manager() }