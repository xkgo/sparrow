package assertpanic

/**
如果断言错误则直接 panic，一般系统启动的时候需要
*/
func NotNil(v interface{}, msg interface{}) {
	if nil == v {
		panic(msg)
	}
}

func NotBlank(v string, msg interface{}) {
	if len(v) < 1 {
		panic(msg)
	}
}
