package main

import (
	"fmt"
	"unsafe"
)

/*
结构体对齐，帮助go 更好地分配内存和 使用cpu读取
*/
type TerraformResource struct {
	Cloud               string // 16 bytes
	Name                string // 16 bytes
	HaveDSL             bool   //  1 byte
	PluginVersion       string // 16 bytes
	IsVersionControlled bool   //  1 byte
	TerraformVersion    string // 16 bytes
	ModuleVersionMajor  int32  //  4 bytes
}

type TerraformResourceOrder struct {
	ModuleVersionMajor  int32  //  4 bytes
	HaveDSL             bool   //  1 byte
	IsVersionControlled bool   //  1 byte
	Cloud               string // 16 bytes
	Name                string // 16 bytes
	PluginVersion       string // 16 bytes
	TerraformVersion    string // 16 bytes

}

func main() {
	var d TerraformResource
	d.Cloud = "aws-singapore" //存储字段 使用的空间与 字段值没有关系
	d.Name = "ec2"
	d.HaveDSL = true
	d.PluginVersion = "3.64"
	d.TerraformVersion = "1.1"
	d.ModuleVersionMajor = 1
	d.IsVersionControlled = true
	fmt.Println("==============================================================")
	fmt.Printf("Total Memory Usage StructType:d %T => [%d]\n", d, unsafe.Sizeof(d))
	fmt.Println("==============================================================")
	fmt.Printf("Cloud Field StructType:d.Cloud %T => [%d]\n", d.Cloud, unsafe.Sizeof(d.Cloud))
	fmt.Printf("Name Field StructType:d.Name %T => [%d]\n", d.Name, unsafe.Sizeof(d.Name))
	fmt.Printf("HaveDSL Field StructType:d.HaveDSL %T => [%d]\n", d.HaveDSL, unsafe.Sizeof(d.HaveDSL))
	fmt.Printf("PluginVersion Field StructType:d.PluginVersion %T => [%d]\n", d.PluginVersion, unsafe.Sizeof(d.PluginVersion))
	fmt.Printf("ModuleVersionMajor Field StructType:d.IsVersionControlled %T => [%d]\n", d.IsVersionControlled, unsafe.Sizeof(d.IsVersionControlled))
	fmt.Printf("TerraformVersion Field StructType:d.TerraformVersion %T => [%d]\n", d.TerraformVersion, unsafe.Sizeof(d.TerraformVersion))
	fmt.Printf("ModuleVersionMajor Field StructType:d.ModuleVersionMajor %T => [%d]\n", d.ModuleVersionMajor, unsafe.Sizeof(d.ModuleVersionMajor))

	var te = TerraformResourceOrder{}
	te.Cloud = "aws-singapore" //存储字段 使用的空间与 字段值没有关系
	te.Name = "ec2"
	te.PluginVersion = "3.64"
	te.TerraformVersion = "1.1"
	te.ModuleVersionMajor = 1
	te.IsVersionControlled = true
	te.HaveDSL = true

	fmt.Println("==============================================================")
	fmt.Printf("Total Memory Usage StructType:d %T => [%d]\n", te, unsafe.Sizeof(te))
	fmt.Println("==============================================================")
	fmt.Printf("Cloud Field StructType:d.Cloud %T => [%d]\n", te.Cloud, unsafe.Sizeof(te.Cloud))
	fmt.Printf("Name Field StructType:d.Name %T => [%d]\n", te.Name, unsafe.Sizeof(te.Name))
	fmt.Printf("PluginVersion Field StructType:d.PluginVersion %T => [%d]\n", te.PluginVersion, unsafe.Sizeof(te.PluginVersion))
	fmt.Printf("TerraformVersion Field StructType:d.TerraformVersion %T => [%d]\n", te.TerraformVersion, unsafe.Sizeof(te.TerraformVersion))

	fmt.Printf("ModuleVersionMajor Field StructType:d.ModuleVersionMajor %T => [%d]\n", te.ModuleVersionMajor, unsafe.Sizeof(te.ModuleVersionMajor))
	fmt.Printf("ModuleVersionMajor Field StructType:d.IsVersionControlled %T => [%d]\n", te.IsVersionControlled, unsafe.Sizeof(te.IsVersionControlled))
	fmt.Printf("HaveDSL Field StructType:d.HaveDSL %T => [%d]\n", te.HaveDSL, unsafe.Sizeof(te.HaveDSL))

	te2 := te
	te2.Cloud = "ali2"
	fmt.Printf("te2:%v\n", te2)
	fmt.Printf("te:%v\n", te)

	te3 := &te ///改变 te3 将同时改变 te,te3 指向了 te的地址
	te3.Cloud = "HWCloud2"
	/*
		(*te3).Cloud
		te3 为指针类型，取其值的时候必须用括号，因为. 的运算优先级 比 * 高
	*/
	// *te3.Cloud 错误

	fmt.Println("(*te3).Cloud:", (*te3).Cloud, "*te3.Cloud", te3.Cloud)
	fmt.Printf("Cloud:%v te3:%p\n", (*te3).Cloud, te3)
	fmt.Printf("Cloud:%v order:%v te:%v, addr:%p\n", te.Cloud, (te).Cloud, te, &te)

}
