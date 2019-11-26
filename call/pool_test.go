package call

import (
	"fmt"
	"midmsg/model"
	"testing"
	"time"
)

func TestTimerCallPool(t *testing.T) {

}


func TestPutPoolAsyncReturn(t *testing.T) {
	test := model.AsyncReturnInfo{
		ClientIP:"123456",
	}
	for i:=0;i<100;i++{
		go AsyncReturn.PutPoolAsyncReturn(test)
	}
	time.Sleep(2*time.Second)
	//AsyncReturn.Unlock()
	fmt.Println(AsyncReturn.AsyncReturnPool)
	//AsyncReturn.Lock()
}