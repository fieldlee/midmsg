package call

import (
	"fmt"
	"github.com/pborman/uuid"
	"testing"
)

func TestRun(t *testing.T) {
	Run()
}

func TestUid(t *testing.T){
	for i := 0 ; i < 100 ; i++{
		fmt.Println(uuid.New())
	}
	for i := 0 ; i < 100 ; i++{
		fmt.Println(uuid.New())
	}
}