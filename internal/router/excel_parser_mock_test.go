package router

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i github.com/hablof/merchant-experience/internal/router.ExcelParser -o ./internal\router\excel_parser_mock_test.go -n ExcelParserMock

import (
	"io"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/hablof/merchant-experience/internal/models"
)

// ExcelParserMock implements ExcelParser
type ExcelParserMock struct {
	t minimock.Tester

	funcParseProducts          func(r io.Reader) (productUpdates []models.ProductUpdate, productErrs []error, err error)
	inspectFuncParseProducts   func(r io.Reader)
	afterParseProductsCounter  uint64
	beforeParseProductsCounter uint64
	ParseProductsMock          mExcelParserMockParseProducts
}

// NewExcelParserMock returns a mock for ExcelParser
func NewExcelParserMock(t minimock.Tester) *ExcelParserMock {
	m := &ExcelParserMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ParseProductsMock = mExcelParserMockParseProducts{mock: m}
	m.ParseProductsMock.callArgs = []*ExcelParserMockParseProductsParams{}

	return m
}

type mExcelParserMockParseProducts struct {
	mock               *ExcelParserMock
	defaultExpectation *ExcelParserMockParseProductsExpectation
	expectations       []*ExcelParserMockParseProductsExpectation

	callArgs []*ExcelParserMockParseProductsParams
	mutex    sync.RWMutex
}

// ExcelParserMockParseProductsExpectation specifies expectation struct of the ExcelParser.ParseProducts
type ExcelParserMockParseProductsExpectation struct {
	mock    *ExcelParserMock
	params  *ExcelParserMockParseProductsParams
	results *ExcelParserMockParseProductsResults
	Counter uint64
}

// ExcelParserMockParseProductsParams contains parameters of the ExcelParser.ParseProducts
type ExcelParserMockParseProductsParams struct {
	r io.Reader
}

// ExcelParserMockParseProductsResults contains results of the ExcelParser.ParseProducts
type ExcelParserMockParseProductsResults struct {
	productUpdates []models.ProductUpdate
	productErrs    []error
	err            error
}

// Expect sets up expected params for ExcelParser.ParseProducts
func (mmParseProducts *mExcelParserMockParseProducts) Expect(r io.Reader) *mExcelParserMockParseProducts {
	if mmParseProducts.mock.funcParseProducts != nil {
		mmParseProducts.mock.t.Fatalf("ExcelParserMock.ParseProducts mock is already set by Set")
	}

	if mmParseProducts.defaultExpectation == nil {
		mmParseProducts.defaultExpectation = &ExcelParserMockParseProductsExpectation{}
	}

	mmParseProducts.defaultExpectation.params = &ExcelParserMockParseProductsParams{r}
	for _, e := range mmParseProducts.expectations {
		if minimock.Equal(e.params, mmParseProducts.defaultExpectation.params) {
			mmParseProducts.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmParseProducts.defaultExpectation.params)
		}
	}

	return mmParseProducts
}

// Inspect accepts an inspector function that has same arguments as the ExcelParser.ParseProducts
func (mmParseProducts *mExcelParserMockParseProducts) Inspect(f func(r io.Reader)) *mExcelParserMockParseProducts {
	if mmParseProducts.mock.inspectFuncParseProducts != nil {
		mmParseProducts.mock.t.Fatalf("Inspect function is already set for ExcelParserMock.ParseProducts")
	}

	mmParseProducts.mock.inspectFuncParseProducts = f

	return mmParseProducts
}

// Return sets up results that will be returned by ExcelParser.ParseProducts
func (mmParseProducts *mExcelParserMockParseProducts) Return(productUpdates []models.ProductUpdate, productErrs []error, err error) *ExcelParserMock {
	if mmParseProducts.mock.funcParseProducts != nil {
		mmParseProducts.mock.t.Fatalf("ExcelParserMock.ParseProducts mock is already set by Set")
	}

	if mmParseProducts.defaultExpectation == nil {
		mmParseProducts.defaultExpectation = &ExcelParserMockParseProductsExpectation{mock: mmParseProducts.mock}
	}
	mmParseProducts.defaultExpectation.results = &ExcelParserMockParseProductsResults{productUpdates, productErrs, err}
	return mmParseProducts.mock
}

// Set uses given function f to mock the ExcelParser.ParseProducts method
func (mmParseProducts *mExcelParserMockParseProducts) Set(f func(r io.Reader) (productUpdates []models.ProductUpdate, productErrs []error, err error)) *ExcelParserMock {
	if mmParseProducts.defaultExpectation != nil {
		mmParseProducts.mock.t.Fatalf("Default expectation is already set for the ExcelParser.ParseProducts method")
	}

	if len(mmParseProducts.expectations) > 0 {
		mmParseProducts.mock.t.Fatalf("Some expectations are already set for the ExcelParser.ParseProducts method")
	}

	mmParseProducts.mock.funcParseProducts = f
	return mmParseProducts.mock
}

// When sets expectation for the ExcelParser.ParseProducts which will trigger the result defined by the following
// Then helper
func (mmParseProducts *mExcelParserMockParseProducts) When(r io.Reader) *ExcelParserMockParseProductsExpectation {
	if mmParseProducts.mock.funcParseProducts != nil {
		mmParseProducts.mock.t.Fatalf("ExcelParserMock.ParseProducts mock is already set by Set")
	}

	expectation := &ExcelParserMockParseProductsExpectation{
		mock:   mmParseProducts.mock,
		params: &ExcelParserMockParseProductsParams{r},
	}
	mmParseProducts.expectations = append(mmParseProducts.expectations, expectation)
	return expectation
}

// Then sets up ExcelParser.ParseProducts return parameters for the expectation previously defined by the When method
func (e *ExcelParserMockParseProductsExpectation) Then(productUpdates []models.ProductUpdate, productErrs []error, err error) *ExcelParserMock {
	e.results = &ExcelParserMockParseProductsResults{productUpdates, productErrs, err}
	return e.mock
}

// ParseProducts implements ExcelParser
func (mmParseProducts *ExcelParserMock) ParseProducts(r io.Reader) (productUpdates []models.ProductUpdate, productErrs []error, err error) {
	mm_atomic.AddUint64(&mmParseProducts.beforeParseProductsCounter, 1)
	defer mm_atomic.AddUint64(&mmParseProducts.afterParseProductsCounter, 1)

	if mmParseProducts.inspectFuncParseProducts != nil {
		mmParseProducts.inspectFuncParseProducts(r)
	}

	mm_params := &ExcelParserMockParseProductsParams{r}

	// Record call args
	mmParseProducts.ParseProductsMock.mutex.Lock()
	mmParseProducts.ParseProductsMock.callArgs = append(mmParseProducts.ParseProductsMock.callArgs, mm_params)
	mmParseProducts.ParseProductsMock.mutex.Unlock()

	for _, e := range mmParseProducts.ParseProductsMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.productUpdates, e.results.productErrs, e.results.err
		}
	}

	if mmParseProducts.ParseProductsMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmParseProducts.ParseProductsMock.defaultExpectation.Counter, 1)
		mm_want := mmParseProducts.ParseProductsMock.defaultExpectation.params
		mm_got := ExcelParserMockParseProductsParams{r}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmParseProducts.t.Errorf("ExcelParserMock.ParseProducts got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmParseProducts.ParseProductsMock.defaultExpectation.results
		if mm_results == nil {
			mmParseProducts.t.Fatal("No results are set for the ExcelParserMock.ParseProducts")
		}
		return (*mm_results).productUpdates, (*mm_results).productErrs, (*mm_results).err
	}
	if mmParseProducts.funcParseProducts != nil {
		return mmParseProducts.funcParseProducts(r)
	}
	mmParseProducts.t.Fatalf("Unexpected call to ExcelParserMock.ParseProducts. %v", r)
	return
}

// ParseProductsAfterCounter returns a count of finished ExcelParserMock.ParseProducts invocations
func (mmParseProducts *ExcelParserMock) ParseProductsAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmParseProducts.afterParseProductsCounter)
}

// ParseProductsBeforeCounter returns a count of ExcelParserMock.ParseProducts invocations
func (mmParseProducts *ExcelParserMock) ParseProductsBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmParseProducts.beforeParseProductsCounter)
}

// Calls returns a list of arguments used in each call to ExcelParserMock.ParseProducts.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmParseProducts *mExcelParserMockParseProducts) Calls() []*ExcelParserMockParseProductsParams {
	mmParseProducts.mutex.RLock()

	argCopy := make([]*ExcelParserMockParseProductsParams, len(mmParseProducts.callArgs))
	copy(argCopy, mmParseProducts.callArgs)

	mmParseProducts.mutex.RUnlock()

	return argCopy
}

// MinimockParseProductsDone returns true if the count of the ParseProducts invocations corresponds
// the number of defined expectations
func (m *ExcelParserMock) MinimockParseProductsDone() bool {
	for _, e := range m.ParseProductsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ParseProductsMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterParseProductsCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcParseProducts != nil && mm_atomic.LoadUint64(&m.afterParseProductsCounter) < 1 {
		return false
	}
	return true
}

// MinimockParseProductsInspect logs each unmet expectation
func (m *ExcelParserMock) MinimockParseProductsInspect() {
	for _, e := range m.ParseProductsMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ExcelParserMock.ParseProducts with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ParseProductsMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterParseProductsCounter) < 1 {
		if m.ParseProductsMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ExcelParserMock.ParseProducts")
		} else {
			m.t.Errorf("Expected call to ExcelParserMock.ParseProducts with params: %#v", *m.ParseProductsMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcParseProducts != nil && mm_atomic.LoadUint64(&m.afterParseProductsCounter) < 1 {
		m.t.Error("Expected call to ExcelParserMock.ParseProducts")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ExcelParserMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockParseProductsInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ExcelParserMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ExcelParserMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockParseProductsDone()
}
