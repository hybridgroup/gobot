package digispark

//#include "littleWire_util.h"
import "C"

func Delay(duration uint) {
	C.delay(C.uint(duration))
}
