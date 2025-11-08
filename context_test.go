package belajargolangcontext

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	bg := context.Background()
	fmt.Println(bg)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")
	contextG := context.WithValue(contextF, "g", "G")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)

	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextF.Value("b"))

	fmt.Println(contextA.Value("b"))
}

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func ()  {
		defer close(destination)
		counter := 1
		for {
			select{
			case <- ctx.Done():
				return 
			default:
			destination <- counter
			counter++
			time.Sleep(1 * time.Second)
			}
		}
	}()
	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination{
		fmt.Println("Counter", n)
		if n == 10{
			break
		}
	}
	cancel()

	fmt.Println(runtime.NumGoroutine())
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Gourutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5 * time.Second)
	defer cancel()

	destination := CreateCounter(ctx)
	fmt.Println("Total Gourutine", runtime.NumGoroutine())
	for n := range destination{
		fmt.Println("Count", n)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Total Gourutine", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Gourutine", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5 * time.Second))
	defer cancel()

	destination := CreateCounter(ctx)
	fmt.Println("Total Gourutine", runtime.NumGoroutine())
	for n := range destination{
		fmt.Println("Count", n)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Total Gourutine", runtime.NumGoroutine())
}

func TestSyncGroup(t *testing.T){
	group := sync.WaitGroup{}

	for i := 0; i < 10; i++{
		group.Add(1)
		go func(num int) {
			defer group.Done()
			fmt.Println("Gorotin ke-", num)
		}(i)
	}

	group.Wait()
	fmt.Println("Selesai")
}

func TestMutex(t *testing.T){
	saldo := 0
	var mutex sync.Mutex
	var group sync.WaitGroup

	for i :=0; i < 100; i++{
		group.Add(1)
		go func(num int) {
			defer group.Done()
			mutex.Lock()
			saldo++
			mutex.Unlock()
		}(i)
	}

	group.Wait()
	fmt.Println("Saldo Akhir", saldo)
}

// 