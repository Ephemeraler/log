package thlog

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func TestProduction(t *testing.T) {
	logger, err := NewProduction("test.log", 1, 1, 1)
	if err != nil {
		t.Errorf("NewProduction() error = %v", err)
		return
	}
	defer logger.Close()

	// logger, err := NewDevelopment()
	// if err != nil {
	// 	t.Errorf("NewProduction() error = %v", err)
	// 	return
	// }
	// defer logger.Close()

	for i := 0; i < 8; i++ {
		logger.zap.With(
			zap.String("url", fmt.Sprintf("www.test%d.com", i)),
			zap.String("name", "jimmmyr"),
			zap.Int("age", 23),
			zap.String("agradege", "no111-000222"),
		).Info("test info ")
	}
}
