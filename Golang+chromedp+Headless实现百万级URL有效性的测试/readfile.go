package main


import (
    "fmt"
    "github.com/360EntSecGroup-Skylar/excelize"
)

func loadfile(txtpath string) []string {
    // f, err := excelize.OpenFile("F:\\GO_code\\test_xlsx\\top-2m.xlsx")
    f, err := excelize.OpenFile(txtpath)
    if err != nil {
        fmt.Println(err)
    }
    
    rows, err := f.GetRows("top-1m")
    var str []string
    // fmt.Printf(string(rows[0][0]))
    for _, row := range rows {
        s := "http://" + row[1]
        str = append(str, s)
    }
    return str
}