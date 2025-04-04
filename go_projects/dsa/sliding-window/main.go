package main

import (
	"fmt"
	"math"
	"strings"
)

//finding the longest substring with no more than k distinct characters
func findLongestSubstring(value string, distinctChar int) string{

	storageMap := make(map[string]int) 
	arr := strings.Split(value, "")
	fmt.Println(arr)
	windowStart := 0
	longestSubString := math.MinInt32 
	var keySubstring string

	for i := 0; i < len(arr); i++ {
		//increment the counter of the char in the substring
		storageMap[arr[i]] += 1

		//check if the no of unique char in the map surpasses the distinctChar
		for(len(storageMap) >= distinctChar){
			//loop over the values in the map to find the longest substring
			//{"A": 3, "H": 2, "I": 3, "B": 1, "C": 1}
			for key, value := range storageMap {
				if(value >= longestSubString) {
					longestSubString = value
					keySubstring = key
				}
				//remove the char from the map when the its counter is zero
				if storageMap[key] == 0{
					delete(storageMap, key)
				}			
			}

			//slide the windowStart by deducting the current char counter
			//  & increasing the windowStart
			storageMap[arr[windowStart]] -= 1
			windowStart++
		}
	}

	return keySubstring
}

func findRepeatedDnaSequences(s string) []string {
	seen := make(map[string]bool)
	repeated := make(map[string]bool)
	var result []string

	for i := 0; i <= len(s)-10; i++ {
		substr := s[i : i+10]
		if seen[substr] {
			repeated[substr] = true
		} else {
			seen[substr] = true
		}
	}

	for seq := range repeated {
		result = append(result, seq)
	}

	return result
}

func main(){
	substring := "AAAHHIIIBC"

	fmt.Printf("Longest substring is %s", findLongestSubstring(substring, 3))

	//DNA SEQUENCES
	s1 := "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"
	fmt.Printf("Dna sequence: %v", findRepeatedDnaSequences(s1))

	s2 := "AAAAAAAAAAAAA"
	fmt.Printf("Dna sequence: %v", findRepeatedDnaSequences(s2))
}