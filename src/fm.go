package main

import (
  "os"
  "fmt"
  "sort"
  "os/exec"
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

func load_panel(active int) {
  dir, err := os.Open(panels[active].path);
  panel, err := dir.Readdir(-1);
  sort.Slice(panel, func(i, j int) bool { return panel[i].Name()[0] != '.' })
  sort.Slice(panel, func(i, j int) bool { return panel[i].IsDir() })
  defer dir.Close(); if err != nil { return }
  load_panel := []map[string]interface{}{}
  load_panel = append(load_panel, map[string]interface{}{
    "name":"..","size":"","date":"","dir": true, "exec":""})
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
      "inserted": 0,
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
  lpl := len(panels[LEFT].path); lpr := len(panels[RIGHT].path); ovl := ""; ovr := ""
  if len(panels[LEFT].path) > int(COLS/2)-2 { lpl = int(COLS/2)-5; ovl = "~" }
  if len(panels[RIGHT].path) > int(COLS/2)-2 { lpr = int(COLS/2)-5; ovr = "~" }
  print_message(2, 0, termbox.ColorBlue | termbox.AttrBold, termbox.ColorDefault, panels[LEFT].path[:lpl] + ovl)
  print_message(2+int(COLS/2), 0, termbox.ColorBlue | termbox.AttrBold, termbox.ColorDefault, panels[RIGHT].path[:lpr] + ovr)
  for row := 0; row < ROWS-2; row++ {
    buffer_row := row + row_offset
    if buffer_row < len(panel) {
      if row < ROWS-2 {
        fgcolor := termbox.ColorWhite;
        if panel[buffer_row]["dir"] == true { fgcolor = termbox.ColorYellow | termbox.AttrBold
        } else if panel[buffer_row]["exec"] == true { fgcolor = termbox.ColorGreen | termbox.AttrBold
        };if panel[buffer_row]["inserted"] == 1 { fgcolor = termbox.ColorRed | termbox.AttrBold }
        overflow := ""; limit := len(panel[buffer_row]["name"].(string));
        if limit > 12 { limit = 12; overflow = "~" }
        name := panel[buffer_row]["name"].(string); size := ""
        if panel[buffer_row]["dir"] == true { if row > 0 { size = "FOLDER" }
        } else { size = panel[buffer_row]["size"].(string)
          if row > 0 { size += " Bytes" }
        };date := panel[buffer_row]["date"].(string)
        if active == LEFT {
          print_message(1, row+1, fgcolor, termbox.ColorDefault, name[:limit]+overflow)
          print_message(int(COLS/3)-len(size)+1, row+1, termbox.ColorCyan, termbox.ColorDefault, size)
          print_message(int(COLS/2)-len(date)-1, row+1, termbox.ColorCyan, termbox.ColorDefault, date)
        } else {
          print_message(offset, row+1, fgcolor, termbox.ColorDefault, panel[buffer_row]["name"].(string)[:limit]+overflow)
          print_message(int(COLS/3)-len(size)+offset, row+1, termbox.ColorCyan, termbox.ColorDefault, size)
          print_message(int(COLS/2)-len(date)-2+offset, row+1, termbox.ColorCyan, termbox.ColorDefault, date)
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

func execute_command() { ROWS--
  termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
  display_borders()
  display_files()
  print_message(1, ROWS, termbox.ColorDefault, termbox.ColorDefault, "NEW DIRECTORY: ")
  termbox.SetCursor(16, ROWS)
  termbox.Flush()
  dir_name := ""
  dir_name_loop:
  for {
    event := get_event()
    switch event.Key {
      case termbox.KeyEsc: break dir_name_loop
      case termbox.KeyEnter:
        if len(dir_name) > 0 {
          cmd := exec.Command("mkdir", panels[active_panel].path + split_ch + dir_name)
          err := cmd.Run(); if err != nil { termbox.SetCursor(-1,-1); return }
          load_panel(LEFT); load_panel(RIGHT)
        };break dir_name_loop
      case termbox.KeySpace: dir_name += " "
      case termbox.KeyBackspace: if len(dir_name) > 0 { dir_name = dir_name[:len(dir_name)-1] }
      case termbox.KeyBackspace2: if len(dir_name) > 0 { dir_name = dir_name[:len(dir_name)-1] }
    };if event.Ch != 0 {
      dir_name += string(event.Ch)
      print_message(16, ROWS, termbox.ColorWhite, termbox.ColorDefault, dir_name)
    };
    dir_name_length := 0
    for _,ch := range dir_name { if ch > 0 { dir_name_length++} }
    termbox.SetCursor(16 + dir_name_length, ROWS)
    for i := 16 + len(dir_name); i < COLS; i++ { termbox.SetChar(i, ROWS, rune(' ')) }
    termbox.Flush()
  };ROWS++; termbox.SetCursor(-1, -1)
}

func process_keypress() {
  event := get_event()
  current_row := &panels[active_panel].current_row
  current_panel := panels[active_panel].panel
  switch event.Key {
    case termbox.KeyEsc: {
      termbox.Close();
      os.Exit(0)
    }
    case termbox.KeyArrowDown: if *current_row < len(current_panel)-1 { *current_row++ }
    case termbox.KeyArrowUp: if *current_row > 0 { *current_row-- }
    case termbox.KeyArrowRight: active_panel = RIGHT
    case termbox.KeyArrowLeft: active_panel = LEFT
    case termbox.KeyEnter: {
      if current_panel[*current_row]["name"] == ".." {
        split_path := strings.Split(panels[active_panel].path, split_ch)
        panels[active_panel].path = strings.Join(split_path[:len(split_path)-1], split_ch)
        if panels[active_panel].path == "" { panels[active_panel].path = split_ch}
        load_panel(active_panel)
      } else if current_panel[*current_row]["dir"] == true {
        is_root := split_ch
        if panels[active_panel].path == is_root { is_root = "" }
        panels[active_panel].path += is_root + current_panel[*current_row]["name"].(string)
        load_panel(active_panel); *current_row = 0;
      }
    }
    case termbox.KeyInsert: {
      if *current_row == 0 { *current_row++; break }
      flip_insert := panels[active_panel].panel[*current_row]["inserted"].(int) ^ 1
      panels[active_panel].panel[*current_row]["inserted"] = flip_insert
      if *current_row < len(current_panel)-1 { *current_row++ }
    }
    case termbox.KeyHome: {
      for _, entry := range panels[active_panel].panel {
        if entry["inserted"] == 1 {
          src_path := panels[active_panel].path + split_ch
          dst_path := panels[active_panel^1].path + split_ch
          src_entry := src_path + entry["name"].(string)
          dst_entry := dst_path + entry["name"].(string)
          cmd := exec.Command("cp", "-r", src_entry, dst_entry)
          err := cmd.Run(); if err != nil { break }
          load_panel(LEFT); load_panel(RIGHT)
        }
      }
    }
    case termbox.KeyDelete: {
      for _, entry := range panels[active_panel].panel {
        if entry["inserted"] == 1 {
          src_path := panels[active_panel].path + split_ch
          src_entry := src_path + entry["name"].(string)
          cmd := exec.Command("rm", "-f", "-r", src_entry)
          err := cmd.Run(); if err != nil { break }
          load_panel(LEFT); load_panel(RIGHT)
        }
      }; *current_row = 0
    }
    case termbox.KeyEnd: { execute_command() }
    case termbox.KeyPgup: *current_row = 0
    case termbox.KeyPgdn: *current_row = len(current_panel)-1
  }
}

func run_file_manager() {
  err := termbox.Init()
  if err != nil { fmt.Println(err); os.Exit(1) }
  panels[LEFT].path, err = os.Getwd()
  panels[RIGHT].path, err = os.Getwd()
  if err != nil { termbox.Close(); fmt.Println("Error reading CWD"); os.Exit(1) }
  load_panel(LEFT)
  load_panel(RIGHT)
  if runtime.GOOS == "windows" { split_ch = "\\" }
  for {
    COLS, ROWS = termbox.Size()
    if COLS < 78 { COLS = 78 }
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
    display_borders()
    display_files()
    termbox.Flush()
    process_keypress()
  }
}

func main() { run_file_manager() }