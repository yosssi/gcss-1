package gcss

// github.com/yosssi/gcss binding for glub.
// No Configuration required.

import (
	"bytes"
	"sync"

	"github.com/omeid/slurp/s"
	"github.com/yosssi/gcss"
)

func Compile(c *s.C) s.Job {
	return func(in <-chan s.File, out chan<- s.File) {

		//Because programs block, zip is not an streaming archive, we don't want to block.
		var wg sync.WaitGroup
		for file := range in {

			buff := new(bytes.Buffer)
			content := file.Content

			go func(file s.File) {
				wg.Add(1)
				defer wg.Done()
				n, err := gcss.Compile(buff, content)
				if err != nil {
					c.Println(err)
				}

				stat := s.FileInfoFrom(file.Stat)
				stat.SetSize(int64(n))

				file.Content = buff
				file.Stat = stat
				out <- file
			}(file)

		}

		wg.Wait()
	}
}
