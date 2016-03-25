package helpers

import ( 
    "time"
    "fmt"
)

func DaysSinceNow(t time.Time) int {
    return int(time.Since(t).Hours())/24
}

func FormatDate(t time.Time) string{
    return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
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