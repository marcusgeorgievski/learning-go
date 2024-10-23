package main

import "fmt"

// 1

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func makePerson(firstName, lastName string, age int) Person {
	return Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}
func makePersonPointer(firstName, lastName string, age int) *Person {
	return &Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

// 2

func UpdateSlice(ss []string, s string) {
	ss[len(ss)-1] = s
	fmt.Println("UpdateSlice:", ss)
}
func GrowSlice(ss []string, s string) {
	ss = append(ss, s)
	fmt.Println("GrowSlice:", ss)
}

// 3

func makePeople() []Person {
	const LIMIT = 10_000_000
	people := make([]Person, 0, LIMIT)

	for i := 0; i < LIMIT; i++ {
		newPerson := Person{
			FirstName: "M",
			LastName:  "G",
			Age:       22,
		}
		people = append(people, newPerson)
	}

	return people
}

// main

func main() {
	// {
	// 	p := makePerson("Marcus", "Georgievski", 22)
	// 	px := makePersonPointer("Marcus", "Georgievski", 22)

	// 	fmt.Println(p)
	// 	fmt.Println(px)
	// }

	// {
	// 	ss := make([]string, 0, 3)
	// 	ss = append(ss, "a", "b", "c")
	// 	fmt.Println("main:", ss)
	// 	UpdateSlice(ss, "updated")
	// 	fmt.Println("main:", ss)
	// 	GrowSlice(ss, "updated")
	// 	fmt.Println("main:", ss)
	// }

	{
		// go build GOGC=x CODEBUG=gctrace=1
		// change capacity of people from 0 to LIMIT, play with GOCG
		people := makePeople()
		fmt.Println(len(people))
	}
}
