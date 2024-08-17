package main

import (
	"container/heap"
	"fmt"
	"strings"
)

// Car represents a car with its registration number and color
type Car struct {
	Registration string
	Color        string
}

// Carpark represents the parking lot
type Carpark struct {
	Slots      map[int]*Car     // Map to store cars by slot number
	EmptySlots IntHeap          // Min-heap for available slots
	MaxSlots   int              // Maximum number of slots
	NextSlot   int              // Next slot number to use if heap is empty
	ColorMap   map[string][]int // Map to store slots by color
	RegMap     map[string]int   // Map to store slot number by registration number
}

// IntHeap implements heap.Interface for a min-heap of integers
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// CreateParkingLot initializes the parking lot with the given number of slots
func (cp *Carpark) CreateParkingLot(n int) {
	cp.Slots = make(map[int]*Car)
	cp.EmptySlots = make(IntHeap, 0, n)
	cp.ColorMap = make(map[string][]int)
	cp.RegMap = make(map[string]int)
	cp.MaxSlots = n
	cp.NextSlot = 1

	for i := 1; i <= n; i++ {
		heap.Push(&cp.EmptySlots, i)
	}
	fmt.Printf("Created a parking lot with %d slots\n", n)
}

// Park parks a car in the parking lot
func (cp *Carpark) Park(registration string, color string) {
	var slotNo int

	if cp.EmptySlots.Len() > 0 {
		slotNo = heap.Pop(&cp.EmptySlots).(int)
	} else if cp.NextSlot <= cp.MaxSlots {
		slotNo = cp.NextSlot
		cp.NextSlot++
	} else {
		fmt.Println("Sorry, parking lot is full")
		return
	}

	if _, exists := cp.Slots[slotNo]; exists {
		fmt.Println("Sorry, parking lot is full")
		return
	}

	cp.Slots[slotNo] = &Car{Registration: registration, Color: color}
	cp.ColorMap[color] = append(cp.ColorMap[color], slotNo)
	cp.RegMap[registration] = slotNo

	fmt.Printf("Allocated slot number: %d\n", slotNo)
}

// Leave frees up a slot
func (cp *Carpark) Leave(slotNo int) {
	if car, exists := cp.Slots[slotNo]; exists {
		delete(cp.Slots, slotNo)
		heap.Push(&cp.EmptySlots, slotNo)

		// Remove slot from ColorMap
		cp.removeSlotFromColorMap(car.Color, slotNo)

		// Remove registration from RegMap
		delete(cp.RegMap, car.Registration)

		fmt.Printf("Slot number %d is free\n", slotNo)
	} else {
		fmt.Println("Slot not found")
	}
}

// removeSlotFromColorMap helper function to remove a slot number from the color map
func (cp *Carpark) removeSlotFromColorMap(color string, slotNo int) {
	colorSlots := cp.ColorMap[color]
	for i, s := range colorSlots {
		if s == slotNo {
			cp.ColorMap[color] = append(colorSlots[:i], colorSlots[i+1:]...)
			if len(cp.ColorMap[color]) == 0 {
				delete(cp.ColorMap, color)
			}
			return
		}
	}
}

// Status prints the current status of the parking lot
func (cp *Carpark) Status() {
	fmt.Println("Slot No. Registration No Colour")
	for i := 1; i <= cp.MaxSlots; i++ {
		if car, ok := cp.Slots[i]; ok {
			fmt.Printf("%d        %s   %s\n", i, car.Registration, car.Color)
		}
	}
}

// RegistrationNumbersForColor returns registration numbers of all cars with a particular color
func (cp *Carpark) RegistrationNumbersForColor(color string) {
	slotNos, exists := cp.ColorMap[color]
	if !exists || len(slotNos) == 0 {
		fmt.Println("Not found")
		return
	}

	regNumbers := make([]string, 0, len(slotNos))
	for _, slotNo := range slotNos {
		if car, exists := cp.Slots[slotNo]; exists {
			regNumbers = append(regNumbers, car.Registration)
		}
	}

	fmt.Println(strings.Join(regNumbers, ", "))
}

// SlotNumbersForColor returns slot numbers of all slots where a car of a particular color is parked
func (cp *Carpark) SlotNumbersForColor(color string) {
	slotNos, exists := cp.ColorMap[color]
	if !exists || len(slotNos) == 0 {
		fmt.Println("Not found")
		return
	}

	slotNosStr := make([]string, 0, len(slotNos))
	for _, slotNo := range slotNos {
		slotNosStr = append(slotNosStr, fmt.Sprintf("%d", slotNo))
	}

	fmt.Println(strings.Join(slotNosStr, ", "))
}

// SlotNumberForRegistrationNumber returns the slot number for a car with a given registration number
func (cp *Carpark) SlotNumberForRegistrationNumber(registration string) {
	slotNo, exists := cp.RegMap[registration]
	if !exists {
		fmt.Println("Not found")
		return
	}

	fmt.Println(slotNo)
}

func main() {
	cp := &Carpark{}
	cp.CreateParkingLot(10)

	cp.Park("KA-01-HH-1234", "White")
	cp.Park("KA-01-HH-9999", "White")
	cp.Park("KA-01-BB-0001", "Black")
	cp.Park("KA-01-HH-7777", "Red")
	cp.Park("KA-01-HH-2701", "Blue")
	cp.Park("KA-01-HH-3141", "Black")
	cp.Leave(4)
	cp.Status()
	cp.Park("KA-01-P-333", "White")
	cp.Park("DL-12-AA-9999", "White")

	cp.RegistrationNumbersForColor("White")
	cp.SlotNumbersForColor("White")
	cp.SlotNumberForRegistrationNumber("KA-01-HH-3141")
	cp.SlotNumberForRegistrationNumber("MH-04-AY-1111")
}
