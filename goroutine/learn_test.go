package belajar_goroutine

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func HelloWeb(number int) {
	fmt.Println("Hello " + strconv.Itoa(number))
}

func TestGoroutine(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		go HelloWeb(i)
	}
	fmt.Println("Main function")
	time.Sleep(2 * time.Second)
}

func giveMessage(channel chan string) {
	time.Sleep(1 * time.Second)
	channel <- "Main function finished"
}

func TestChannel(t *testing.T) {
	channel := make(chan string)
	defer close(channel)

	go giveMessage(channel)

	data := <-channel
	fmt.Println(data)
}

func TestRange(t *testing.T) {
	channel := make(chan string)
	go func() {
		for i := 0; i < 10; i++ {
			channel <- strconv.Itoa(i)
		}
		close(channel)
	}()

	for data := range channel {
		fmt.Println(data)
	}
}

func TestSelect(t *testing.T) {
	chan1 := make(chan string)
	chan2 := make(chan string)

	defer close(chan1)
	defer close(chan2)
	go giveMessage(chan1)
	go giveMessage(chan2)

	count := 0

	for {
		select {
		case data := <-chan1:
			fmt.Println("Received on chan1:", data)
			count++
		case data := <-chan2:
			fmt.Println("Received on chan2:", data)
			count++
		default:
			fmt.Println("Waiting")
		}

		if count >= 2 {
			break
		}
	}

}

func TestRaceCondition(t *testing.T) {
	counter := 0

	for i := 0; i < 1000; i++ {
		go func() {
			counter++
		}()
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Final counter:", counter)
}

func TestMutex(t *testing.T) {
	counter := 0
	mutex := sync.Mutex{}

	for i := 0; i < 1000; i++ {
		go func() {
			mutex.Lock()
			counter++
			mutex.Unlock()
		}()
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Final counter:", counter)
}

type BankAccount struct {
	balance int
	mutex   sync.RWMutex
}

func (b *BankAccount) reduceBalance(amount int) {
	b.mutex.Lock()
	if b.balance <= 0 {
		fmt.Println("Insufficient balance")
		b.mutex.Unlock()
		return
	}
	b.balance -= amount
	b.mutex.Unlock()
}

func (b *BankAccount) getBalance() int {
	b.mutex.RLock()
	balance := b.balance
	b.mutex.RUnlock()
	return balance
}

func TestReadWriteMutex(t *testing.T) {
	account := BankAccount{balance: 50}

	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				account.reduceBalance(1)
				fmt.Printf("balance after reduce %v \n", account.getBalance())
			}
		}()
	}

	time.Sleep(5 * time.Second)
	fmt.Printf("last balance %v \n", account.getBalance())
}

type UserBalance struct {
	Mutex   sync.Mutex
	Name    string
	Balance int
}

func (u *UserBalance) Lock() {
	u.Mutex.Lock()
}

func (u *UserBalance) Unlock() {
	u.Mutex.Unlock()
}

func (u *UserBalance) Change(amount int) {
	u.Balance += amount
}

func Transfer(user1 *UserBalance, user2 *UserBalance, amount int) {
	user1.Lock()
	fmt.Println("Lock user 1", user1.Name)
	user1.Change(-amount)

	time.Sleep(1 * time.Second)

	user2.Lock()
	fmt.Println("Lock user 2", user2.Name)
	user2.Change(amount)

	time.Sleep(1 * time.Second)

	user1.Unlock()
	user2.Unlock()
}

func TestDeadlock(t *testing.T) {
	user1 := UserBalance{Name: "Eko", Balance: 100}
	user2 := UserBalance{Name: "Budi", Balance: 100}

	go Transfer(&user1, &user2, 50)
	go Transfer(&user2, &user1, 50)

	time.Sleep(5 * time.Second)
	// harusnya balance eko dan budi 100 tapi keduanya 50 karena terjadi deadlock di logic user2.lock()
	fmt.Println(user1.Balance, user2.Balance)
}

func RunAsync(group *sync.WaitGroup) {
	group.Add(1)

	fmt.Println("Running asynchronously")

	group.Done()
}

func TestWaitGroup(t *testing.T) {
	group := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		go RunAsync(group)
	}
	group.Wait()
	fmt.Println("All goroutines finished")
}

func OnlyOnce(count *int) {
	*count = *count + 1
}

func TestOnce(t *testing.T) {
	once := sync.Once{}
	group := sync.WaitGroup{}
	counter := 0

	for i := 0; i < 10; i++ {
		group.Add(1)
		go func() {
			once.Do(func() { OnlyOnce(&counter) })
			group.Done()
		}()
	}
	group.Wait()
	fmt.Printf("counter value: %v\n", counter)
}

func TestPool(t *testing.T) {
	group := sync.WaitGroup{}

	//pool without default value
	//pool := sync.Pool{}

	// pool with default value
	// pool sudah aman dari race condition
	pool := sync.Pool{
		New: func() interface{} {
			return ""
		},
	}

	eko := "Eko"
	kur := "Kurniawan"

	pool.Put(&eko)
	pool.Put(&kur)

	for i := 0; i < 3; i++ {
		group.Add(1)
		go func() {
			data := pool.Get()
			fmt.Println(data)
			pool.Put(data)
			group.Done()
		}()
	}

	group.Wait()
}

func addToMap(data *sync.Map, group *sync.WaitGroup, value int) {
	group.Add(1)
	data.Store(value, value)
	group.Done()
}

func TestMap(t *testing.T) {
	// aman dari race condition
	data := sync.Map{}
	group := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		go addToMap(&data, &group, i)
	}

	group.Wait()
	data.Range(func(key, value any) bool {
		fmt.Println(value)
		return true
	})
}

func TestAtomic(t *testing.T) {
	// primitive data aman dari race condition
	var counter int64
	group := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		group.Add(1)
		go func() {
			atomic.AddInt64(&counter, 2)
			group.Done()
		}()
	}

	group.Wait()

	fmt.Println("Final counter:", atomic.LoadInt64(&counter))
}

func TestTimer(t *testing.T) {
	timer := time.NewTimer(5 * time.Second)
	fmt.Printf("time now %v\n", time.Now())

	timeReceived := <-timer.C

	fmt.Println("Timer expired", timeReceived)

	group := sync.WaitGroup{}
	group.Add(1)
	time.AfterFunc(2*time.Second, func() {
		fmt.Println("Timer expired after 2 seconds")
		group.Done()
	})
	group.Wait()
}

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)

	for data := range ticker.C {
		fmt.Println("Tick at", data)
	}
}

func TestGoMax(t *testing.T) {
	totalCpu := runtime.NumCPU()
	fmt.Println(totalCpu)
	totalThread := runtime.GOMAXPROCS(-1)
	fmt.Println(totalThread)
	totalGoroutine := runtime.NumGoroutine()
	fmt.Println(totalGoroutine)
}
