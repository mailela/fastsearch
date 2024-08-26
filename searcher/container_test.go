package searcher

import (
	"fmt"
	"testing"
)

func TestContainer_Init(t *testing.T) {
	c := &Container{
		Dir:   "./data/",
		Debug: true,
	}
	err := c.Init()
	if err != nil {
		panic(err)
	}

	test := c.GetDataBase("icms3")

	fmt.Println(test.GetIndexCount())
	fmt.Println(test.GetDocumentCount())

	all := c.GetDataBases()
	for name, engine := range all {
		fmt.Println(name)
		fmt.Println(engine)
	}
}
