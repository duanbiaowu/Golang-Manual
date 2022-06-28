package spider

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetGoVersion(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockSpider := NewMockSpider(mockCtl)
	mockSpider.EXPECT().GetBody().Return("go1.8.3")
	goVer := GetGoVersion(mockSpider)

	assert.Equal(t, "go1.8.3", goVer)
}
