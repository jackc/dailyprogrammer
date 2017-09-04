package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func main() {
	bufstdout := bufio.NewWriter(os.Stdout)
	defer bufstdout.Flush()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), " ")
		nums := make([]int64, len(words))
		for i, w := range words {
			var err error
			nums[i], err = strconv.ParseInt(w, 10, 64)
			if err != nil {
				log.Fatalf("%w is not a number", w)
			}
		}

		sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })

		answers := [][]int64{}
		for i := 0; i < len(nums); i++ {
			for j := i + 1; j < len(nums); j++ {
				for k := j + 1; k < len(nums); k++ {
					if nums[i]+nums[j]+nums[k] == 0 {
						answers = append(answers, []int64{nums[i], nums[j], nums[k]})
					}
				}
			}
		}

		if len(answers) > 0 {
			sort.Slice(answers, func(i, j int) bool {
				return (answers[i][0] < answers[j][0]) ||
					(answers[i][0] == answers[j][0] && answers[i][1] < answers[j][1]) ||
					(answers[i][0] == answers[j][0] && answers[i][1] == answers[j][1] && answers[i][2] < answers[j][2])
			})

			uniqAnswers := make([]string, 0, len(answers))
			var prevAnswer []int64
			for _, a := range answers {
				if !reflect.DeepEqual(prevAnswer, a) {
					s := fmt.Sprintf("%d %d %d", a[0], a[1], a[2])
					uniqAnswers = append(uniqAnswers, s)
					prevAnswer = a
				}
			}

			for _, a := range uniqAnswers {
				fmt.Fprintln(bufstdout, a)
			}

			fmt.Fprintln(bufstdout)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

// while gets
//   nums = $_.split.map { |s| Integer(s) }.sort

//   answers = []
//   (0...nums.size).each do |i|
//     ((i+1)...nums.size).each do |j|
//       ((j+1)...nums.size).each do |k|
//         answers.push [nums[i], nums[j], nums[k]] if nums[i] + nums[j] + nums[k] == 0
//       end
//     end
//   end

//   answers.sort.uniq.each do |a|
//     puts a.join(" ")
//   end

//   puts
// end
