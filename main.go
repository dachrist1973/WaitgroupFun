package waitgroupfun

import (
	"fmt"
	"sync"
)

func waitfunc(id int, wg *sync.WaitGroup) {
	// Call Done on the wait group. Better to let the thread call done
	// and best to defer to cover your self in case you write code that possibly
	// misses calling done that way you know it will be called.
	defer (*wg).Done()
	fmt.Printf("wait id is %d\n", id)
}

func waitfunc2(id int) {
	fmt.Printf("wait id is %d\n", id)
}

// WaitGroupPreload - The following function add items to a wait group prior loop execution that fires off the
// go routines
func WaitGroupPreload() {

	var wg sync.WaitGroup
	// Add item to wait group before running the for loop \
	// this was seen on an interview question.
	wg.Add(10)

	for i := 0; i < 10; i++ {
		// Call the wait func above passing in the waitgroup.
		// Note wait groups should be passed in as a pointer.
		go waitfunc(i, &wg)

	}

	// Need to wait if you do not then the go routines may never start.
	wg.Wait()
	print("Done\n")

}

// WaitGroupLoopLoad - This function does the waitgroup add in the for loop.
func WaitGroupLoopLoad() {

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go waitfunc(i, &wg)

	}
	wg.Wait()
	print("Done\n")

}

// WaitGroupLoopLoad2
func WaitGroupLoopLoad2() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go waitfunc2(i)
		wg.Done()
	}
	wg.Wait()
	print("Done\n")

}

// Lock - wrapper func to lock a mutex
func Lock(mt *sync.Mutex) {
	(*mt).Lock()
}

// Unlock - wrapper function to unlock a mutex
func Unlock(mt *sync.Mutex) {
	(*mt).Unlock()
}

// mutexLocking - demostrates how to pass mutexes around to be use in other functions.
// mutexes need to be passed in as pointers.
func mutexLocking(mt *sync.Mutex, playmap *map[string]int, name string) {
	Lock(mt)
	(*playmap)[name]++
	Unlock(mt)
}

// MutexFun1 - this shows basic mutex locking of a shared resource
// need to use waitGroups for the go routines
func MutexFun1() {
	var mt sync.Mutex
	// shared resource
	playmap := map[string]int{"a": 0, "b": 0}
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go mutexLocking(&mt, &playmap, "a")
		wg.Done()

	}
	wg.Wait()
	println(playmap["a"])
	println(playmap["b"])
}

// This code demonstrates the use of channels letting two
// go routines pass values to each other.

// The data we pass between the go routines.
type ChanStruct struct {
	Nodnum   int
	Nodename string
}

// ReadChannel - reads the date from the channel.
func ReadChannel(c chan ChanStruct) {
	for {
		node, ok := <-c
		if ok {
			fmt.Printf("Read Node: %d, Name: %s\n", node.Nodnum, node.Nodename)
		} else {
			break
		}

	}

}

// WriteChannel - writes data to the channel takes two channels one channel acts as
// an atomic to show the thread has started.
func WriteChannel(c chan ChanStruct, nodes []ChanStruct, start chan bool) {
	// write that we have started the thread.
	start <- true
	// close the channel form other writes
	close(start)

	// Lets write some data to the channel
	for _, node := range nodes {
		// Close the channel after adding all the data.
		if node.Nodename == "Close" {
			close(c)
		} else {
			// Add the data to channel
			c <- node
			fmt.Printf("Added Node num %d, Node name %s\n", node.Nodnum, node.Nodename)
		}
	}

}

// ChannelFun
func ChannelFun() {
	// Create main channel to pass data into
	StrChan := make(chan ChanStruct, 3)

	// create channel to
	start := make(chan bool, 1)
	nodes := []ChanStruct{
		{
			Nodnum:   1,
			Nodename: "test1",
		},
		{
			Nodnum:   2,
			Nodename: "test2",
		},
		{
			Nodnum:   3,
			Nodename: "test3",
		},
		{
			Nodnum:   4,
			Nodename: "Close",
		},
	}

	go WriteChannel(StrChan, nodes, start)
	for {
		started := <-start
		if started == true {
			break
		}
	}
	go ReadChannel(StrChan)

}
