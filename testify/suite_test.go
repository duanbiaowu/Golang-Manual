package testify

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MyTestSuit struct {
	suite.Suite
	count uint32
}

type MyHttpSuite struct {
	suite.Suite
	recorder *httptest.ResponseRecorder
	mux      *http.ServeMux
}

func (s *MyHttpSuite) SetupSuite() {
	s.recorder = httptest.NewRecorder()
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", index)
	s.mux.HandleFunc("/greeting", greeting)
}

func (s *MyTestSuit) SetupSuite() {
	fmt.Println("SetupSuite")
}

func (s *MyTestSuit) TearDownSuite() {
	fmt.Println("TearDownSuite")
}

func (s *MyTestSuit) SetupTest() {
	fmt.Printf("SetupTest test count:%d\n", s.count)
}

func (s *MyTestSuit) TearDownTest() {
	s.count++
	fmt.Printf("TearDownTest test count:%d\n", s.count)
}

func (s *MyTestSuit) BeforeTest(suiteName, testName string) {
	fmt.Printf("BeforeTest suite:%s test:%s\n", suiteName, testName)
}

func (s *MyTestSuit) AfterTest(suiteName, testName string) {
	fmt.Printf("AfterTest suite:%s test:%s\n", suiteName, testName)
}

func (s *MyTestSuit) TestExample() {
	fmt.Println("TestExample")
}

func TestExample(t *testing.T) {
	suite.Run(t, new(MyTestSuit))
}

func (s *MyHttpSuite) TestIndex() {
	request, _ := http.NewRequest("GET", "/", nil)
	s.mux.ServeHTTP(s.recorder, request)

	s.Assert().Equal(s.recorder.Code, 200, "get index error")
	s.Assert().Contains(s.recorder.Body.String(), "Hello World", "body error")
}

func (s *MyHttpSuite) TestGreeting() {
	request, _ := http.NewRequest("GET", "/greeting", nil)
	request.URL.RawQuery = "name=dj"

	s.mux.ServeHTTP(s.recorder, request)

	s.Assert().Equal(s.recorder.Code, 200, "greeting error")
	s.Assert().Contains(s.recorder.Body.String(), "welcome, dj", "body error")
}

func TestHTTP(t *testing.T) {
	suite.Run(t, new(MyHttpSuite))
}
