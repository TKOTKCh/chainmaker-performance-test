package main

//go:wasmexport sum
func sum(x, y int64) int64 {
	return x + y
}
func main() {

}
