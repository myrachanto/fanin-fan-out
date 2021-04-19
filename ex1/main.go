package main

import (
	"fmt"
	"sync"
)
type result struct{
	name string
	weeks int
	days int
} 
type data struct {
	name string
	days int
}
func main() {
	m := map[string]int { "Robert" : 30, "John" : 475, "Elly" : 1022, "Bob" : 99 }
	//fan out part
	res := Getdata(m)
	worker1 := worker(res)
	worker2 := worker(res)
	worker3 := worker(res)
	//fanin part
	ch, workers := fanin(worker1,worker2,worker3)
	for c := range ch {
		fmt.Printf("%v has being working for the company for %v weeks and %v day(s)\n", c.name, c.weeks, c.days)
	}
	fmt.Printf("%v workers were employed to do the job", workers)


}
func fanin(chans ...<-chan result) (<-chan result, int){
	out := make(chan result)
	wg := &sync.WaitGroup{}
	wg.Add(len(chans))
	workers := len(chans)
	for _, c := range chans {
		go func(res <-chan result){
			for r := range res{
				out<-r
			}
			wg.Done()
		}(c)
	}
	go func(){
		wg.Wait()
		close(out)
	}()
	return out, workers
}
func Getdata(m map[string]int)  <-chan data{
	out := make(chan data)
	go func(){
		for k,v := range m {
			res := data{
					name: k,
					days: v,
			}
			out<-res
		}
		close(out)
	}()
	return out
 
}
func worker(data <-chan data)  <-chan result{
	out := make(chan result)
	go func(){
		for c := range data {
			res := result{
					name: c.name,
					weeks: c.days/7,
					days: c.days%7,
			}
			out<-res
		}
		close(out)
	}()
	return out
}