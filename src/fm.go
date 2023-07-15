package main

import (
  "os"
  "fmt"
  "runtime"
  "strings"
  "strconv"
  "github.com/nsf/termbox-go"
  "github.com/mattn/go-runewidth"
)

var COLS, ROWS int
var LEFT = 0
var RIGHT = 1
var active_panel = LEFT

type Panel struct {
  path string
  panel []map[string]interface{}
  current_row, row_offset int
}

var panels = []Panel {
  {
    path: "/home/cmk/Desktop",
    panel: []map[string]interface{}{},
    current_row: 0,
    row_offset: 0,
  },
  {
    path: "/bin",
    panel: []map[string]interface{}{},
    current_row: 0,
    row_offset: 0,
  },
}

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


func load_panel(active int) {
  dir, err := os.Open(panels[active].path);
  panel, err := dir.Readdir(-1);
  defer dir.Close(); if err != nil { return }
  load_panel := []map[string]interface{}{}
  load_panel = append(load_panel, map[string]interface{}{"name": "..", "dir": true})
  for row := 0; row < len(panel); row++ {
    exec := false
    if runtime.GOOS == "windows" {
      exec = strings.Contains(panel[row].Name(), ".exe")
    } else if runtime.GOOS == "linux" { exec = panel[row].Mode()&0111 != 0 }
    info := map[string]interface{}{
      "name": panel[row].Name(),
      "size": strconv.FormatInt(panel[row].Size(), 10),
      "date": panel[row].ModTime().Format("2006-01-02"),
      "dir": panel[row].IsDir(),
      "exec": exec,
      "inserted": false,
    }; load_panel = append(load_panel, info)
  };panels[active].panel = load_panel
}

func display_panel(offset, active int) {
  var panel = []map[string]interface{}{}
  row_offset := 0; current_row := 0
  col_from := 0; col_to := 0
  panel = panels[active].panel
  row_offset = panels[active].row_offset
  current_row = panels[active].current_row
  if current_row < row_offset { row_offset = current_row }
  if current_row >= row_offset + ROWS-2 { row_offset = current_row - ROWS+2+1 }
  panels[active].row_offset = row_offset
  print_message(2, 0, termbox.ColorBlue | termbox.AttrBold, termbox.ColorDefault, panels[LEFT].path)
  print_message(2+int(COLS/2), 0, termbox.ColorBlue | termbox.AttrBold, termbox.ColorDefault, panels[RIGHT].path)
  for row := 0; row < ROWS-2; row++ {
    buffer_row := row + row_offset
    if buffer_row < len(panel) {
      if row < ROWS-2 {
        fgcolor := termbox.ColorWhite;
        if panel[buffer_row]["dir"] == true { fgcolor = termbox.ColorYellow | termbox.AttrBold
        } else if panel[buffer_row]["exec"] == true { fgcolor = termbox.ColorGreen | termbox.AttrBold }
        if active == LEFT {
          print_message(1, row+1, fgcolor, termbox.ColorDefault, panel[buffer_row]["name"].(string))
        } else {
          print_message(1+offset, row+1, fgcolor, termbox.ColorDefault, panel[buffer_row]["name"].(string))
        }
      }
    }
    if row == current_row - row_offset && active_panel == active {
      _, ROWS := termbox.Size();
      if row >= ROWS { continue }
      if active_panel == LEFT { col_from = 1; col_to = int(COLS/2)-1
      } else { col_from = int(COLS/2)+1; col_to = COLS-1 }
      for col := col_from; col < col_to; col++ {
        current_cell := termbox.GetCell(col, row+1)
        termbox.SetCell(col, row+1, current_cell.Ch, current_cell.Fg, termbox.ColorBlue)
      }
    }
  }
}

func display_files() {
  display_panel(0, LEFT)
  display_panel(int(COLS/2)+1, RIGHT)
}

func print_message(x, y int, fg, bg termbox.Attribute, message string) {
  for _, c := range message {
    termbox.SetCell(x, y, c, fg, bg)
    x += runewidth.RuneWidth(c)
  }
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
  current_row := &panels[active_panel].current_row
  current_panel := panels[active_panel].panel
  switch event.Key {
    case termbox.KeyEsc: termbox.Close(); os.Exit(0)
    case termbox.KeyArrowDown: if *current_row < len(current_panel)-1 { *current_row++ }
    case termbox.KeyArrowUp: if *current_row > 0 { *current_row-- }
    case termbox.KeyArrowRight: active_panel = RIGHT
    case termbox.KeyArrowLeft: active_panel = LEFT
    case termbox.KeyEnter: {
      split_ch := "/"
      if current_panel[*current_row]["name"] == ".." {
        split_path := strings.Split(panels[active_panel].path, split_ch)
        panels[active_panel].path = strings.Join(split_path[:len(split_path)-1], split_ch)
        if panels[active_panel].path == "" { panels[active_panel].path = "/"}
        load_panel(active_panel)
      } else if current_panel[*current_row]["dir"] == true {
        panels[active_panel].path += split_ch + current_panel[*current_row]["name"].(string)
        load_panel(active_panel); *current_row = 0
      }
    }
  }
}

func run_file_manager() {
  err := termbox.Init()
  if err != nil { fmt.Println(err); os.Exit(1) }
  load_panel(LEFT)
  load_panel(RIGHT)
  for {
    COLS, ROWS = termbox.Size()
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
    display_panels()
    display_files()
    termbox.Flush()
    process_keypress()
  }
}

func main() { run_file_manager() }