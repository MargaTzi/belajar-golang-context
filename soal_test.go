package belajargolangcontext

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

type User struct {
	sync.RWMutex
	Name string
}

//write
func(u *User) Tambah(alamat string){
	u.RWMutex.Lock()
	u.Name += alamat
	u.RWMutex.Unlock()
}

//read
func (u *User) Baca()string{
	u.RWMutex.RLock()
	defer u.RWMutex.RUnlock()
	return u.Name
}

func TestRW(t *testing.T) {
	user := &User{}
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 3 * time.Second)
	defer cancel()

	var group sync.WaitGroup
	//Menulis
	group.Add(1)
	go func() {
		defer group.Done()
		for {
			select{
			case <-ctx.Done():
				fmt.Println("Penulis berhenti")
				return
			default:
				user.Tambah("Dito")
				fmt.Println("Penulis: Menambahkan Nama")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	
	//Membaca
	for i:=0; i<3;i++{
		group.Add(1)
		go func(id int) {
			defer group.Done()
			for{
				select{
				case <- ctx.Done():
					fmt.Println("Pembaca", id, "berhenti")
					return
				default:
					name := user.Baca()
					fmt.Println("Pembaca", id, "melihat", name)
					time.Sleep(500 * time.Millisecond)
				}
			}
		}(i)
	}

	group.Wait()
	fmt.Println("Selesai")
}

func GenerateNumbers(limit int) chan int{
	num := make(chan int)

	go func() {
		defer close(num)
		for i:= 1; i<=limit; i++{
			num <- i
			// time.Sleep(1 * time.Second)
		}
	}()

	return num
}

func TestGenerateNumbers(t *testing.T){
	data := GenerateNumbers(100000)

	for n := range data{
		fmt.Println("Menerima Angka ke -", n)
	}
	fmt.Println("selesai")
}

func Producer(limit int) chan int{
	channel := make(chan int)
	go func() {
		defer close(channel)
		for i := 1; i<=limit; i++{
			channel <- i
		}
	}()
	return channel
}

func Consumer(channel <-chan int){
	for num := range channel{
		fmt.Println("Menerima ke-", num)
	}
}

func TestChannelbasic(t *testing.T) {
	var wg sync.WaitGroup
	data := Producer(10)
	
	wg.Go(func() {
		Consumer(data)
	})

	wg.Wait()
}