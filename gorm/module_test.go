package gorm

import (
	"fmt"
	"github.com/delfimarime/fx/config"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"testing"
)

func TestNew(t *testing.T) {
	c := config.Config{
		Gorm: map[string]config.GormDatasource{
			"default": {
				Type: SqlLiteType,
				URL:  fmt.Sprintf("%s/default.db", t.TempDir()),
			},
			"test": {
				Type: SqlLiteType,
				URL:  fmt.Sprintf("%s/test.db", t.TempDir()),
			},
		},
	}
	app := fx.New(New(c), fx.Invoke(func(seq []*gorm.DB) {
		for _, each := range seq {
			fmt.Println(each)
		}
	}))

	app.Run()
}
