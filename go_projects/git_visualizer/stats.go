package main

import (
	"fmt"
	"sort"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)
const outofRange = 99999
const daysInLastSixMonths = 183
const weeksInLastSixMonths = 26
type column []int
//Generates the stats for a given git status repository

//1. get the lists of commits
//2. given the commits, generate the graph

//processRepositories given a user email, returns the commits made in the last 6 months
func processRepositories(email string) map[int]int {
	filepath := getDotFilePath()
	repos := parseFileLinesToSlice(filepath)
	daysInMap := daysInLastSixMonths

	commits := make(map[int]int, daysInMap)
	for i := daysInMap; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(email, path, commits)
	}

	return commits
}

//getBeginningOfDay given a time.Time calculates the start time of that day
func getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay :=  time.Date(year, month, day, 0,0,0,0, t.Location())
	return startOfDay
} 

func countDaysSinceDate(date time.Time) int {
	days := 0
	now := getBeginningOfDay(time.Now())

	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++

		if days > daysInLastSixMonths {
			return outofRange
		}
	}
	
	return days
}
//calOffset determines and returns the amount of days missing to fill the last row of the stats graph
//used to calculate the right commit place of a commit in our commits map
func calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}

	return offset
}

//fillCommits given a repository found in path, gets the commits and 
//puts them in the commits map, returning it when completed
func fillCommits(email string, path string, commits map[int]int) map[int]int {
	//instantiate a git repo object from path
	repo,err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	//get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}

	//get the commits history starting from the HEAD
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}

	//iterate the commits
	offset := calcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSinceDate(c.Author.When) + offset

	  //acts like an authenfication check.
		//@todo: extract this feature to a more robust function 
		if c.Author.Email != email {
			return nil
		}

		if daysAgo != outofRange {
			commits[daysAgo]++
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return commits
}

//printCommitsStats prints the commits stats
//1. Sort the map
//2. Generate the columns
//3. Print each column
func printCommitStats(commits map[int]int) {
	keys := sortMapIntoSlice(commits)
	cols := buildCols(keys, commits)

	printCells(cols)
}

//takes a map and returns a slice with the map keys ordered by their integer value
func sortMapIntoSlice(m map[int]int) []int {
	//order map
	//Store the keys in slice in sorted order
	var keys []int 
	for k := range m {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	return keys
}

//Generate columns
//generates a map with rows and columns ready to be printed to screen
func buildCols(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := int(k / 7)
		dayinweek := k % 7

		if dayinweek == 0 {
			col = column{}
		} 

		col = append(col, commits[k])

		if dayinweek == 6 {
			cols[week] = col
		}
	}

	return cols
}

//prints the month names in the first line, determining when the month changed between switching weeks
func printMonths(){
	week := getBeginningOfDay(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("     ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("   ")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}

	fmt.Printf("\n")
}
// printDayCol given the day number (0 is sunday) prints the day name
func printDayCol(day int){
	out := "    "
	switch day {
	case 1: 
		out = "Mon"
	case 2: 
		out = "Wed"
	case 3:
		out = "Fri"
	}

	fmt.Printf(out)
}

func printCell(val int, today bool) {
	escape := "\033[0;37;30m"

	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m"
	case val >= 10:
		escape = "\033[1;30;45m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := " %d"
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

//printcells prints the cells of the graph
func printCells(cols map[int]column) {
	printMonths()
	for j := 6; j >= 0; j-- {
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {
			if i == weeksInLastSixMonths + 1 {
				printDayCol(j)
			}
			
			if col, ok := cols[i]; ok {
				//special case today

				if i == 0 && j == calcOffset() -1 {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false)
						continue
					}
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}

}

func stats(email string) {
	commits := processRepositories(email)
	printCommitStats(commits)
}
