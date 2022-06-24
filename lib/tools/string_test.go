package tools

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString(t *testing.T) {
	Convey("string to byte", t, func() {
		str := "test"
		byteStr := StringToBytes(str)
		So(byteStr, ShouldResemble, []byte(str))
		Convey("byte to String", func() {
			So(BytesToString(byteStr), ShouldEqual, str)
		})
	})
}
