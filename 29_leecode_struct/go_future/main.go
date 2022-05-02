//传值 接收对象并返回 对象
func MakeNewFoo(f Foo) (Foo, error) {
	f.Field1 = 'New value'
	f.Field2 = f.Field2 + 1
	return f, nil
	}
//通过引用传递,  这将接收一个指向Foo的指针并改变原始对象
func MutateFoo(f *Foo( error {
	f.Field1 = "New val"
	f.Field2 = 2
	return nil
	}




