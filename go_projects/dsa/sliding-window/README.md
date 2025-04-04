Sliding window technique is a problem-solving technique that works by maintaining a ==window==, a ==contiguous part== of the data that respects the constraints. Then we keep moving that window by extending it and shrinking it while respecting constraints until we finish the whole input.

```go
//Example 1: Find the max price of 3 books in a shelf

//[Brute Force Way]:
func bestTotalPrice(prices []int, k int){
	//here we store the maxtotal captured while looping
	//every five books category
	maxtotal = 0

	//at this point we are looping in this manner
	//[1, 2, 3, 4, 5, 6, 7, 8, 9]
	//the category to be calc will be as follows:
	// 1. {1, 2, 3}
	// 2. {2, 3, 4}
	// 3. {3, 4, 5}
	// 4. {4, 5, 6}
	// 5. {5, 6, 7}
	// 6. {6, 7, 8}
	// 7. {7, 8, 9}
	for i, price := range len(prices)- k + 1 { // O(n - k + 1)
		//here we find the sum of the elem in the group
		total = sum(prices[i:i+k]) //time complexity = O(n)
		maxtotal = max(maxtotal
```


### Characteristics of Sliding Window Pattern

1.  Things we iterate over are sequential, contiguous(items that are grouped together in a subset), example: strings, arrays, linked list
2. Some of the constraints are: min, max, longest, shortest, contained


## Questions Variants

#### 1. Fixed Length
- Max sum ==subarray== of size k

>[!info] Example: Max SubArray Technique
>``` go
>   func maxSubArray(intArr []int, group int){
> 	 //calculate the total of the intial subset of window
> 	 //& capture it as the maxTotal for the initialization
> 	  total := sum(intArr[i : i + group])
> 	  maxTotal := total
> 	 for i,num := range (len(intArr) - group + 1){
> 		 //sliding window procedure:
> 		 //1. Deduct the first elem within the window
> 		 //2. Add the adjacent elem after the last elem in the window 
> 		 total -= intArr[i]
> 		 total += intArr[i + group]
> 		 //compare against the maxtotal
> 		 maxTotal = max(maxTotal, total)
> 	 }
> 	 return maxTotal
>   }
>   func sum(elem []int) int {
> 	  total := 0
> 	  for _, num := range elem {
> 		  total  += num
> 	  }
> 	  return total
>   }


![[Sliding-Window(Static).excalidraw]]
### 2. Dynamic Variant
- Smallest sum >= to some value s

>[!info] Example: Find the smallest subarray which is >= to the target sum(5)
>```go
>    func smallestSubArray(targetSum int, arr []int){
> 	   currentSum := 0
> 	   windowStart := 0
> 	   minWindowSize := math.maxInteger
> 	   for windowEnd, num := range arr {
> 		   //keep track of the current sum of the elem in the window
> 		   currentSum += num
> 		   //whenever we find that the current sum >= the target sum we wanna decrease the starting point of the window
> 		  for(currentSum >= targetSum){
> 			  minWindowSize = math.min(minWindowSize, windowEnd - windowStart + 1)
> 			  currentSum -= arr[windowStart]
> 			  windowStart++
> 		  }
> 	   }
> 	  return minWindowSize
>    }
>
>```




![[DynaminSlidingWindow.excalidraw]]

### 3. Dynamic variant w/ Auxillary Data Structure

- Longest substring with no more than k distinct characters

```go
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
			// & increasing the windowStart
			storageMap[arr[windowStart]] -= 1
			windowStart++

		}

	}
		return keySubstring
}
```

![[DynaminSlidingWindow.excalidraw]]



### Applications of Sliding Window Technique

### 🔍 **1. Network Monitoring & Packet Analysis**

- **Use case**: Detecting anomalies or calculating average bandwidth usage.
    
- **Example**: Monitor the number of packets sent over a network every 5 seconds. Instead of recalculating from scratch, the sliding window helps maintain an updated total as new packets arrive and old ones expire.

---

### 🎵 **2. Audio Signal Processing**

- **Use case**: Noise reduction, speech recognition, or pitch detection.
- **Example**: Audio is processed in small overlapping chunks (windows) to apply filters or extract features like MFCCs for machine learning models.
---

### 📹 **3. Video Stream Compression**

- **Use case**: Frame difference detection.
- **Example**: To compress video, compare the current frame to the last _n_ frames using a sliding window to store recent history and detect changes.


### 🧪 **4. DNA/Genome Analysis**

- **Use case**: Searching for specific gene patterns.
- **Example**: A sliding window of length _k_ is used to find repeating or unique k-mers (substrings of length _k_) in DNA sequences.

---

### 🧮 **5. Image Processing / Computer Vision**

- **Use case**: Object detection or feature extraction.
- **Example**: Apply filters (like edge detection) across an image by sliding a 3x3 or 5x5 window across all pixels.

