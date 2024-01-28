package ReportAnnualLeaveService

import (
	"fmt"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"sort"
	"sync"
	"time"
)

func (input reportAnnualLeaveService) worker(id, year int, month <-chan int, results chan<- []repository.EmployeeLeaveReportModel, err chan<- errorModel.ErrorModel, ctrl chan struct{}, wg *sync.WaitGroup) {
	var (
		day  = 1
		zero = 0
	)

	defer wg.Done()
	for itemMonth := range month {
		select {
		case <-ctrl:
			fmt.Printf("Worker [%d] canceled\n", id)
			return
		default:
			fmt.Printf("Worker [%d] processing for [month %d] [year %d]\n", id, itemMonth, year)
			timeNow := time.Date(year, time.Month(itemMonth), day, zero, zero, zero, zero, time.UTC)

			//--- Main Process
			datas, errTemp := input.mainProcessReportAnnualLeave(timeNow)
			if errTemp.Error != nil {
				fmt.Printf("Error in Worker [%d], [month %d] [year %d], canceling all workers\n", id, itemMonth, year)
				err <- errTemp
				close(ctrl)
				return
			}

			results <- datas
		}
	}
}

func (input reportAnnualLeaveService) doGetReportAnnualLeaveInYear(year int) (resultData []repository.EmployeeLeaveReportModel, err errorModel.ErrorModel) {
	var (
		wg          sync.WaitGroup
		monthOnYear = 12
		jobs        = make(chan int, monthOnYear)
		results     = make(chan []repository.EmployeeLeaveReportModel, monthOnYear)
		errs        = make(chan errorModel.ErrorModel, 1)
		ctrl        = make(chan struct{})
	)

	//--- [1] Run Routine
	for i := 1; i <= monthOnYear; i++ {
		wg.Add(1)
		go input.worker(i, year, jobs, results, errs, ctrl, &wg)
	}

	//--- [2] Month Looping
	for j := 1; j <= monthOnYear; j++ {
		jobs <- j
	}

	close(jobs)

	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()

	//--- [3] Collect Return
	for itemErr := range errs {
		err = itemErr
		return
	}

	//--- [4] Collect Results
	for result := range results {
		resultData = append(resultData, result...)
	}

	//--- [5] Sort By Name
	sort.Slice(resultData, func(i, j int) bool {
		return resultData[i].RowNumber.Int64 < resultData[j].RowNumber.Int64
	})

	err = errorModel.GenerateNonErrorModel()
	return
}
