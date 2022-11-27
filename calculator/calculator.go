package calculator

import (
	"log"
	"sync"

	"finHubPipeline/structs"

	"database/sql"
)

type MovingAverage struct {
	slotsFilled bool
	Window      int
	valPos      int
	values      []float64
	db          *sql.DB
	mux         sync.RWMutex
}

func New(window int, db *sql.DB) *MovingAverage {
	return &MovingAverage{
		Window:      window,
		db:          db,
		values:      make([]float64, window),
		valPos:      0,
		slotsFilled: false,
	}
}

func (ma *MovingAverage) CalculateData(m structs.MsgData) <-chan float64 {
	ma.mux.Lock()
	defer ma.mux.Unlock()
	var wg sync.WaitGroup
	oc := make(chan float64)
	wg.Add(1)
	go func() {
		ma.values[ma.valPos] = m.P
		ma.valPos = (ma.valPos + 1) % ma.Window
		if !ma.slotsFilled && ma.valPos == 0 {
			ma.slotsFilled = true
		}
		wg.Done()
		oc <- ma.Avg(m.S)
		close(oc)
	}()
	wg.Wait()
	return oc
}

func (ma *MovingAverage) Avg(symbol string) float64 {
	ma.mux.RLock()
	defer ma.mux.RUnlock()
	var sum = float64(0)
	values := ma.filledValues()
	if values == nil {
		return 0
	}
	n := len(values)

	for _, value := range values {
		sum += value
	}
	avg := sum / float64(n)

	if ma.db != nil {
		ma.saveDataToPgSQL(avg, symbol)
	}

	return avg
}

func (ma *MovingAverage) filledValues() []float64 {
	var c = ma.Window - 1

	if !ma.slotsFilled {
		c = ma.valPos - 1
		if c < 0 {
			return nil
		}
	}
	return ma.values[0 : c+1]
}

func (ma *MovingAverage) saveDataToPgSQL(avg float64, symbol string) {
	sqlStatement := `
			INSERT INTO logs (avg_value, symbol)
			VALUES ($1, $2)`
	_, err := ma.db.Exec(sqlStatement, avg, symbol)
	if err != nil {
		log.Fatalf("Error exec sql query, %v", err)
	}
}
