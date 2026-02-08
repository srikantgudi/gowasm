# gowasm

## This is a WASM app created in Golang

## Developing WASM app in Golang is pretty straightforward.

### Declare functions in the `main`

```golang
- js.Global().Set("setZoneTime", js.FuncOf(setZoneTime))
- js.Global().Set("getZoneTime", js.FuncOf(getZoneTime))
- js.Global().Set("setLocalTime", js.FuncOf(setLocalTime))
- js.Global().Set("wctof", js.FuncOf(wctof))
- js.Global().Set("wftoc", js.FuncOf(wftoc))
- js.Global().Set("wGetReportData", js.FuncOf(wGetReportData))
```

### Define them next, for example

```golang
func wctof(this js.Value, args []js.Value) interface{} {
  val, _ := strconv.ParseFloat(args[0].String(), 32)
  resultId := args[1].String()
  totemp := (val * 1.8) + 32.0
  getElement(resultId).Set("innerHTML", fmt.Sprintf("%.2f C = %.2f F", val, totemp))
  return nil
}
```

#### That's all the coding

### Build the app

```bash
GOOS=js GOARCH=wasm go build -o gowasm.wasm -ldflags "-w -s"
```

### Run the app using any http server, I used this

```
npx serve .
```