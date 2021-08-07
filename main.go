package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func main() {
	CheckArgs("url")
	url := os.Args[1]
	Info("git clone %s  --recursive", url)
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
	})

	CheckIfError(err)

	ref, err := r.Head()
	CheckIfError(err)

	iter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)
	err = iter.ForEach(func(c *object.Commit) error {
		author := c.Author.Name
		when := c.Author.When
		fs, e := c.Stats()
		CheckIfError(e)
		for _, stat := range fs {
			addToStats(author, when.Year(), int(when.Month()), when.Day(), stat.Addition, stat.Deletion)
		}
		return nil
	})
	CheckIfError(err)
	showStats()
}

type Stat struct {
	Addition int
	Deletion int
	Files    int
}

type Key struct {
	Author string
	Year   int
	Month  int
	Day    int
}

var statistics map[Key]Stat = make(map[Key]Stat)

func addToStats(author string, year, month, day, add, del int) {
	key := Key{author, year, month, day}
	value := statistics[key]
	value.Addition += add
	value.Deletion += del
	value.Files++
	statistics[key] = value
}

func showStats() {
	active_days := make(map[string]int)
	files_per_author := make(map[string]int)
	adds := make(map[string]int)
	dels := make(map[string]int)

	for key := range statistics {
		value := statistics[key]
		active_days[key.Author]++
		files_per_author[key.Author] += value.Files
		adds[key.Author] += value.Addition
		dels[key.Author] += value.Deletion
	}

	keys := make([]string, 0, len(active_days))
	for k := range active_days {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Printf("%-20s %-6s %-10s %-10s %-10s\n", "Author", "Days", "hocs", "adds", "dels")
	for _, k := range keys {
		n := active_days[k]
		if n > 0 {
			a := adds[k]
			d := dels[k]
			hocs := (a + d) / active_days[k]
			author := strings.TrimSpace(k)
			fmt.Printf("[%-20.20s]|%6d %10d %10d %10d\n", author, active_days[k], hocs, a, d)
		}

	}
}
