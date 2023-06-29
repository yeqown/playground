package v2

import (
	"sync"
	"time"

	"github.com/playground/golang/gormcs"
)

func (g *gorm2TestSuite) Test_bench_tx() {
	wg := sync.WaitGroup{}
	routine := func() {
		defer wg.Done()
		tx := g.db().Model(&gormcs.UserModel{})
		tx.Begin()

		out := new(gormcs.UserModel)
		if err := tx.Where("id = ?", 1).First(out).Error; err != nil {
			return
		}

		time.Sleep(time.Millisecond * 10)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go routine()
	}

	wg.Wait()
	time.Sleep(10 * time.Second)
}
