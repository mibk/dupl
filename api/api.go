package api

import (
	"sort"

	"github.com/mibk/dupl/internal/clones"
	"github.com/mibk/dupl/internal/job"
	"github.com/mibk/dupl/internal/syntax"
)

// Issue represents two identical parts of files.
type Issue struct {
	From, To Clone
}

// Clone represents duplicated unit of code in a file.
type Clone struct {
	Filename  string
	LineStart int
	LineEnd   int
	Fragment  []byte
}

// Run finds duplicates content of list of files.
func Run(files []string, threshold int) ([]Issue, error) {
	fchan := make(chan string, 1024)
	go func() {
		for _, f := range files {
			fchan <- f
		}
		close(fchan)
	}()
	schan := job.Parse(fchan)
	t, data, done := job.BuildTree(schan)
	<-done

	// finish stream
	t.Update(&syntax.Node{Type: -1})

	mchan := t.FindDuplOver(threshold)
	duplChan := make(chan syntax.Match)
	go func() {
		for m := range mchan {
			match := syntax.FindSyntaxUnits(*data, m, threshold)
			if len(match.Frags) > 0 {
				duplChan <- match
			}
		}
		close(duplChan)
	}()

	return proccessDuplicates(duplChan)
}

func proccessDuplicates(duplChan <-chan syntax.Match) ([]Issue, error) {
	groups := make(map[string][][]*syntax.Node)
	for dupl := range duplChan {
		groups[dupl.Hash] = append(groups[dupl.Hash], dupl.Frags...)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var issues []Issue
	for _, k := range keys {
		uniq := syntax.Unique(groups[k])
		if len(uniq) > 1 {
			duplicates, err := clones.CreateClones(uniq)
			if err != nil {
				return nil, err
			}
			i, err := createIssues(duplicates)
			if err != nil {
				return nil, err
			}
			issues = append(issues, i...)
		}
	}

	return issues, nil
}

func createIssues(duplicates []clones.Clone) ([]Issue, error) {
	sort.Sort(clones.ByNameAndLine(duplicates))
	var issues []Issue
	for i, cl := range duplicates {
		nextCl := duplicates[(i+1)%len(duplicates)]
		issues = append(issues, Issue{
			From: Clone(cl),
			To:   Clone(nextCl),
		})
	}
	return issues, nil
}
