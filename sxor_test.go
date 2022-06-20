package sxor




import (
	"testing"
	"fmt"
	"reflect"
)

func TestReconst(t *testing.T) {
	x, _ := New(2, 2)
	data := make([]byte, 100)
	for i := range data {
		data[i] = byte(i)
	}
	vects_fixed, _ := x.Split(data)
	vects, _ := x.Split(data)
	x.Encode(vects)
	x.Encode(vects_fixed)
	
	
	test_id :=0
	vects[0] = make([]byte, 0)
	x.Reconstruct(vects)
	if !reflect.DeepEqual(vects, vects_fixed){
		fmt.Printf("Failed test_id = %d", test_id)
	}

	test_id =1
	vects[1] = make([]byte, 0)
	x.Reconstruct(vects)
	if !reflect.DeepEqual(vects, vects_fixed){
		fmt.Printf("Failed test_id = %d", test_id)
	}

	test_id =2
	vects[2] = make([]byte, 0)
	x.Reconstruct(vects)
	if !reflect.DeepEqual(vects, vects_fixed){
		fmt.Printf("Failed test_id = %d", test_id)
	}

	test_id =3
	vects[3] = make([]byte, 0)
	x.Reconstruct(vects)
	if !reflect.DeepEqual(vects, vects_fixed){
		fmt.Printf("Failed test_id = %d", test_id)
	}

	test_id =4
	vects[0] = make([]byte, 0)
	vects[1] = make([]byte, 0)
	x.Reconstruct(vects)
	if !reflect.DeepEqual(vects, vects_fixed){
		fmt.Printf("Failed test_id = %d", test_id)
	}
	
	/*
	for i := range vects{
		fmt.Printf("vect[%d]=", i)
		for j := range vects[i]{
			fmt.Printf("%d,",vects[i][j])
		}
		fmt.Print("\n")
	}*/

}
