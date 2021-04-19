package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)
type City struct {
	Name string
	Population int
}
func main(){
	f, err := os.Open("ex2/cities.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	rows := genRows(f)
	smallcity := MinPopulation(60000000)
	//fan out pattern....................../////////
	// 	   __ worker1  
	//rows/__worker2
	//	  \__worker3
	//     __workern...
	worker1 := Uppercase(smallcity(rows))
	worker2 := Uppercase(smallcity(rows))
	worker3 := Uppercase(smallcity(rows))
	worker4 := Uppercase(smallcity(rows))
	//////////////////////////////////////////
	//fan in pattern
	//worker1__
	//worker2__\data
	//worker3__/
	channs, workers := FanIn(worker1,worker2,worker3,worker4)
	for c := range channs {
		fmt.Println(c )
	}
	fmt.Printf("with the %v workers", workers)

}
func Uppercase(cities <-chan City) <-chan City{
	out := make(chan City)
	go func (){
		for c := range cities {
			out<- City{ Name: strings.ToUpper(c.Name), Population: c.Population}
		}
		close(out)
	}()

	return out
}
func MinPopulation(min int) func(<-chan City) <-chan City{
	return func(cities <-chan City) <-chan City{
		out := make(chan City)
		go func (){
			for c := range cities {
				if c.Population < min{
					out<- City{ Name: strings.ToUpper(c.Name), Population: c.Population}
				}
			}
			close(out)
		}()

		return out
	}
}
func Bigcity(cities <-chan City) <-chan City{
	out := make(chan City)
	go func (){
		for c := range cities {
			if c.Population > 40000{
				out<- City{ Name: c.Name, Population: c.Population}
			}
		}
		close(out)
	}()

	return out
}

func FanIn(chans ...<-chan City) (<-chan City, int) {
	out := make(chan City)
	wg := &sync.WaitGroup{}
	wg.Add(len(chans))
	workers := len(chans)
	for _, c := range chans {
		go func(city <-chan City) {
			for r := range city {
				out <- r
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out, workers
}
func genRows(r io.Reader) chan City {
	out := make(chan City)
	go func(){
		reader := csv.NewReader(r)
		_, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}
		for {
			row,err := reader.Read()
			if err == io.EOF{
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			pop, err := strconv.Atoi(row[9])
			if err != nil {
				continue
			}
			out<- City{Name:row[1], Population: pop}
		}
		close(out)

	}()
	return out
}