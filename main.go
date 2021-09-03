package goWeb

import (
	"fmt"

	"github.com/cjie9759/goWeb/controller"
)

func Demo() {
	fmt.Println(NewApp().Get(&controller.Index{}).SetMiddle(MWLog).Run(":17127"))
}
