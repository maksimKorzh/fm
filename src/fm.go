package main

import "os"
import "fmt"
//import "strings"
import "strconv"
import "github.com/nsf/termbox-go"
import "github.com/mattn/go-runewidth"

var COLS, ROWS int
var left = 0
var right = 1
var left_row_offset, right_row_offset, active_panel int
var left_current_row, right_current_row int
var left_panel = []map[string]interface{}{}
var right_panel = []map[string]interface{}{}
var left_path = "/home/cmk"
var right_path = "/bin"

func display_cell(lt, rt, lb, rb [] int) {
  for row := lt[1]+1; row < lb[1]; row++ {
    termbox.SetChar(lt[0], row, '│')
    termbox.SetChar(rt[0], row, '│')
    for col := lt[0]+1; col < rt[0]; col++ {
      termbox.SetChar(col, rt[1], '─')
      termbox.SetChar(col, rb[1], '─')
  }};termbox.SetChar(lt[0], lt[1], '┌')
  termbox.SetChar(lb[0], lb[1], '└')
  termbox.SetChar(rt[0], rt[1], '┐')
  termbox.SetChar(rb[0], rb[1], '┘')
}

func display_panels() {
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


func load_panel(path string, active int) {
  dir, err := os.Open(path);
  panel, err := dir.Readdir(-1);
  defer dir.Close(); if err != nil { return }
  load_panel := []map[string]interface{}{}
  for row := 0; row < len(panel); row++ {
    info := map[string]interface{}{
      "name": panel[row].Name(),
      "size": strconv.FormatInt(panel[row].Size(), 10),
      "date": panel[row].ModTime().Format("2006-01-02"),
      "dir": panel[row].IsDir(),
      "exec": panel[row].Mode().Perm()&0111,
      "inserted": false,
    }; load_panel = append(load_panel, info)
  }
  if active == left { left_panel = load_panel
  } else { right_panel = load_panel }
}

func display_panel(offset, active int) {
  var panel = []map[string]interface{}{}
  row_offset := 0; current_row := 0
  col_from := 0; col_to := 0
  if active == left {
    panel = left_panel
    row_offset = left_row_offset
    current_row = left_current_row
  } else {
    panel = right_panel
    row_offset = right_row_offset
    current_row = right_current_row
  }



  if current_row < row_offset { row_offset = current_row }
  if current_row >= row_offset + ROWS-2 { row_offset = current_row - ROWS+2+1 }
  if active == left { left_row_offset = row_offset
  } else { right_row_offset = row_offset }

  for row := 0; row < ROWS-2; row++ {
    buffer_row := row + row_offset
    if buffer_row < len(panel) {
      if row < ROWS-3 {
        if active == left {
          print_message(1, row+2, termbox.ColorWhite, termbox.ColorBlue, panel[buffer_row]["name"].(string))
        } else {
          print_message(1+offset, row+2, termbox.ColorWhite, termbox.ColorBlue, panel[buffer_row]["name"].(string))
        }
      }
    }


    if row == current_row - row_offset && active_panel == active {
      _, ROWS := termbox.Size();
      if row >= ROWS { continue }
      if active_panel == left { col_from = 1; col_to = int(COLS/2)-1
      } else { col_from = int(COLS/2)+1; col_to = COLS-1 }
      for col := col_from; col < col_to; col++ {
        current_cell := termbox.GetCell(col, row+1)
        termbox.SetCell(col, row+1, current_cell.Ch, termbox.ColorDefault, termbox.ColorYellow)
      }
    }

  }


}

func display_files() {
  display_panel(0, left)
  display_panel(int(COLS/2)+1, right)
}

func print_message(x, y int, fg, bg termbox.Attribute, message string) {
  for _, c := range message {
    termbox.SetCell(x, y, c, fg, bg)
    x += runewidth.RuneWidth(c)
  }//;termbox.Flush()
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
    case termbox.KeyEsc: termbox.Close(); os.Exit(0)
    case termbox.KeyArrowDown:
      if active_panel == 0 {
        if left_current_row < len(left_panel) { left_current_row++ }
      } else { if right_current_row < len(right_panel) { right_current_row++ }}
    case termbox.KeyArrowUp:
      if active_panel == 0 {
        if left_current_row > 0 { left_current_row-- }
      } else { if right_current_row > 0 { right_current_row-- }}
    case termbox.KeyArrowRight: active_panel = right
    case termbox.KeyArrowLeft: active_panel = left
  }
}

func run_file_manager() {
  err := termbox.Init()
  if err != nil { fmt.Println(err); os.Exit(1) }
  load_panel(left_path, left)
  load_panel(right_path, right)
  for {
    COLS, ROWS = termbox.Size()
    termbox.Clear(termbox.ColorWhite, termbox.ColorBlue)
    display_panels()
    display_files()
    termbox.Flush()
    process_keypress()
  }
}

func main() { run_file_manager() }