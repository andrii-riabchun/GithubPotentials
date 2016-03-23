package helpers

import "time"

func DaysSinceNow(t time.Time) int {
    return int(time.Since(t).Hours())/24
}

func RemoveDuplicates(elements []int) []int {
    encountered := make(map[int]bool)
    var result []int
    for v := range elements {
        if encountered[elements[v]] == false {
        encountered[elements[v]] = true
	    result = append(result, elements[v])
        }
    }
    return result
}