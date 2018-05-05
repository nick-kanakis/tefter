package cmd

import (
	"github.com/nicolasmanic/tefter/model"
	"testing"
)

func TestPrintOverview(t *testing.T) {
	cases := []struct {
		notebooks []*model.Notebook
		deep      bool
	}{
		{
			notebooks: []*model.Notebook{
				&model.Notebook{
					ID:    1,
					Notes: nil,
					Title: "title",
				}, &model.Notebook{
					ID:    2,
					Notes: make(map[int64]*model.Note),
					Title: "title2",
				},
			},
			deep: false,
		}, {
			notebooks: []*model.Notebook{
				&model.Notebook{},
			},
			deep: true,
		}, {
			notebooks: []*model.Notebook{},
			deep:      true,
		}, {
			notebooks: nil,
			deep:      true,
		},
	}

	for _, c := range cases {
		printOverview(c.notebooks, c.deep)
	}
}
