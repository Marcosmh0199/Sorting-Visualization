/*
https://en.wikipedia.org/wiki/Linear_congruential_generator
https://en.wikipedia.org/wiki/Go_(programming_language)
https://stackoverflow.com/questions/47256354/windows-how-to-get-screen-resolution-in-golang
https://github.com/gizak/termui
https://syntaxdb.com/ref/go/user-input
https://books.google.co.cr/books?hl=es&lr=&id=p-vpBwAAQBAJ&oi=fnd&pg=PR7&dq=Factorizations+and+primality+testing,+Springer-Verlag,+New+York,+NY,+1989&ots=b8Aa4JhFoH&sig=-hX6x70gAnBWRF1khf989cFg0Xw#v=onepage&q=Factorizations%20and%20primality%20testing%2C%20Springer-Verlag%2C%20New%20York%2C%20NY%2C%201989&f=false
 */

package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/lxn/win"
	"log"
	"math"
	"math/rand"
	"time"
)

var width int = int(win.GetSystemMetrics(win.SM_CXSCREEN) / 8)
var height int = int(win.GetSystemMetrics(win.SM_CYSCREEN) / 24)
var sortChart widgets.BarChart
var primes []int
const MAX_NUMBER_SIZE = 32
const SLICE_SIZE = 50
const MAX_PRIME = 103
func main() {
	slice := randomSlice(101)
	generatePrimes()
	barChartDriver(slice)
}

func remove(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}

func generatePrimes(){
	var booleans = make([]bool, MAX_PRIME)
	for i := range booleans{
		booleans[i] = true
	}
	for i := 2; i*i <= MAX_PRIME; i++{
		if booleans[i] == true {
			for j := i*i; j <= MAX_PRIME; j += i{
				booleans[j] = false
			}
		}
	}
	for i := 2; i < MAX_PRIME; i++{
		if booleans[i]{
			primes = append(primes, i)
		}
	}
}

/*
creates a N size slice with random numbers based on the linear congruential method using only the seed as param
output: slice with N random integers
 */
func randomSlice(seed int) []float64{
	var slice = make([]float64, SLICE_SIZE)
	var m = int(math.Pow(2, MAX_NUMBER_SIZE))  // modulus
	var a = 22695477                    // multiplier
	var c = 1                           //increment
	for i := 0; i < SLICE_SIZE; i++ {
		seed = (a * seed + c) % m
		slice[i] = float64(seed % MAX_NUMBER_SIZE)
	}
	return slice
}

/*
Aux function to scale a number from 0-31 to 0-255
*/
func numberToRed (number int) uint8 {
	return uint8(number * 255 / (MAX_NUMBER_SIZE-1))
}

/*
Aux function to swap two numbers
*/
func swap (a *float64, b *float64){
	temp := *a
	*a = *b
	*b = temp
}

func updateChart(){
	ui.Render(&sortChart)
	time.Sleep(10 * time.Millisecond)
}

//optimized bubbleSort, generates a graph each loop
func bubbleSort(slice []float64) []float64{
	n := len(slice) - 1
	cont := 0
	for true {
		swapped := false
		for i := 0; i < n; i++{
			if slice[i] > slice[i+1]{
				swap(&slice[i], &slice[i+1])
				swapped = true
				updateChart()
			}
		}
		if !swapped{
			break
		}
		n--
		cont++
	}
	return slice
}

func displayHelp(){
	for i := int(win.GetSystemMetrics(win.SM_CYSCREEN) / 24); i > 0; i-- {
		println("")
	}
	println("______________________________________________ ")
	println("|Comandos:                                   | ")
	println("|1: Ejecutar quicksort sobre el slice actual | ")
	println("|2: Ejecutar bubblesort sobre el slice actual| ")
	println("|3: Regresar el slice a su estado original   | ")
	println("|4: Crear un slice nuevo                     | ")
	println("|5: Salir                                    | ")
	println("---------------------------------------------- ")
}

func barChartDriver(slice []float64) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	initSortChart(slice)
	ui.Render(&sortChart)
	displayHelp()
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			sortChart.BarColors= []ui.Color{ui.StandardColors[rand.Int()%len(ui.StandardColors)]}
			switch e.ID {
			case "1":
				sortChart.Title = "QuickSort"
				sortChart.Data = make([]float64, len(slice))
				copy(sortChart.Data, slice)
				quickSort(sortChart.Data)
			case "2":
				sortChart.Title = "BubbleSort"
				sortChart.Data = make([]float64, len(slice))
				copy(sortChart.Data, slice)
				bubbleSort(sortChart.Data)
			case "3":
				sortChart.Title = "Sort"
				sortChart.Data = make([]float64, len(slice))
				copy(sortChart.Data, slice)
				ui.Render(&sortChart)
			case "4":
				slice = randomSlice(primes[rand.Int() % len(primes)])
				sortChart.Data = make([]float64, len(slice))
				copy(sortChart.Data, slice)
				ui.Render(&sortChart)
			case "5":
				return
			}
		}
	}
}

func initSortChart(slice []float64)  {
	sortChart = *widgets.NewBarChart()
	sortChart.Data = slice
	sortChart.Title = "Sort"
	sortChart.SetRect(0, 0, width, height)
	sortChart.BarWidth = 3
	sortChart.BarGap = 0
	sortChart.BarColors = []ui.Color{ui.ColorRed}
	sortChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
}

func quickSort(slice []float64) []float64 {
	if len(slice) < 2 {
		return slice
	}
	low := 0
	high := len(slice)-1

	pivot := rand.Int() % len(slice)

	swap(&slice[pivot], &slice[high])
	for i := range slice {
		if slice[i] < slice[high] {
			swap(&slice[i], &slice[low])
			updateChart()
			low++
		}
	}
	swap(&slice[low], &slice[high])
	quickSort(slice[:low])
	quickSort(slice[low+1:])
	return slice
}
