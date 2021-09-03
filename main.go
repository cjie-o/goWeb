package goWeb

import (
	"fmt"
	"goWeb/controller"
)

func Demo() {
	fmt.Println(NewApp().Get(&controller.Index{}).SetMiddle(MWLog).Run(":17127"))
}
