package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"syscall/js"
	"time"
	_ "time/tzdata"
)

func main() {
	js.Global().Set("setZoneTime", js.FuncOf(setZoneTime))
	js.Global().Set("getZoneTime", js.FuncOf(getZoneTime))
	js.Global().Set("setLocalTime", js.FuncOf(setLocalTime))
	js.Global().Set("wctof", js.FuncOf(wctof))
	js.Global().Set("wftoc", js.FuncOf(wftoc))
	js.Global().Set("wGetReportData", js.FuncOf(wGetReportData))
	fmt.Println("Go-Wasm App ready...")
	select {}
}

// utility functions

func getElement(id string) js.Value {
	return js.Global().Get("document").Call("getElementById", id)
}

func setText(id string, text any) {
	element := getElement(id)
	if !element.IsUndefined() {
		element.Set("innerText", text)
	} else {
		log.Printf("Element with ID '%s' not found", id)
	}
}

// END: utility functions

func setLocalTime(this js.Value, args []js.Value) interface{} {
	now := time.Now()
	getElement("local-time").Set("innerHTML", now.Format("Mon 02-Jan 15:04:05 am MST"))

	locAngHr := (now.Hour()*30 + now.Minute()/2) - 90
	locAngMi := (now.Minute()*6 + now.Second()/10) - 90
	locAngSe := (now.Second() * 6) - 90
	getElement("locthr").Call("setAttribute", "transform", fmt.Sprintf("rotate(%v)", locAngHr))
	getElement("loctmi").Call("setAttribute", "transform", fmt.Sprintf("rotate(%v)", locAngMi))
	getElement("loctse").Call("setAttribute", "transform", fmt.Sprintf("rotate(%v)", locAngSe))
	return nil
}

func setZoneTime(this js.Value, args []js.Value) interface{} {
	zoneId := args[0].String()

	loc, err := time.LoadLocation(zoneId)
	if err != nil {
		log.Fatalln("err:", err.Error())
	}
	now := time.Now()
	ztime := now.In(loc)
	getElement("zone-name").Set("innerHTML", zoneId)
	getElement("zone-time").Set("innerHTML", ztime.Format("Mon 02-Jan 15:04:05 am MST"))

	angHr := (ztime.Hour()*30 + ztime.Minute()/2) - 90
	angMi := (ztime.Minute()*6 + ztime.Second()/10) - 90
	angSe := (ztime.Second() * 6) - 90
	getElement("thr").Call("setAttribute", "transform", fmt.Sprintf("rotate(%v)", angHr))
	getElement("tmi").Call("setAttribute", "transform", fmt.Sprintf("rotate(%v)", angMi))
	getElement("tse").Call("setAttribute", "transform", fmt.Sprintf("rotate(%v)", angSe))
	return nil
}

func getZoneTime(this js.Value, args []js.Value) interface{} {
	zoneID := args[0].String()
	loc, _ := time.LoadLocation(zoneID)
	now := time.Now().In(loc)
	return js.ValueOf(now.Format(time.RFC1123))
}

func wctof(this js.Value, args []js.Value) interface{} {
	val, _ := strconv.ParseFloat(args[0].String(), 32)
	resultId := args[1].String()
	totemp := (val * 1.8) + 32.0
	getElement(resultId).Set("innerHTML", fmt.Sprintf("%.2f C = %.2f F", val, totemp))
	return nil
}

func wftoc(this js.Value, args []js.Value) interface{} {
	val, _ := strconv.ParseFloat(args[0].String(), 32)
	resultId := args[1].String()
	totemp := (val - 32.0) / 1.8
	getElement(resultId).Set("innerHTML", fmt.Sprintf("%.2f F = %.2f C", val, totemp))
	return nil
}
func wGetReportData(this js.Value, args []js.Value) interface{} {
	path := args[0].String()
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/%v", path)
	go func() {
		getElement("dynamic-result").Set("innerHTML", "----- fetching from: ["+path+"] -----")
		time.Sleep(1 * time.Second)
		resp, err := http.Get(url)
		if err != nil {
			getElement("dynamic-result").Set("innerHTML", err.Error())
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			getElement("dynamic-result").Set("innerHTML", err.Error())
			return
		}
		getElement("dynamic-result").Set("innerHTML", string(body))
	}()
	return nil
}
