package sxor




import (
	"testing"
	"fmt"
	"reflect"
	"time"
	"math/rand"
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


func BenchmarkEncode(t *testing.B) {
	x, _ := New(2, 2)
	data := make([]byte, int(1e7))
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	vects, _ := x.Split(data)

	start := time.Now()
	x.Encode(vects)
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Encode time= %v\n", elapsed)
	fmt.Printf("Encode Speed= %v MB/s\n", float64(len(data))/1e6 / elapsed.Seconds())
}

func BenchmarkDecode0(t *testing.B) {
	x, _ := New(2, 2)
	data := make([]byte, int(1e7))
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	vects, _ := x.Split(data)
	x.Encode(vects)

	vects[0] = make([]byte, 0) //drop the 0-th disk
	start := time.Now()
	x.Reconstruct(vects)
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Decode0 time= %v\n", elapsed)
	fmt.Printf("Decode0 Speed= %v MB/s\n", float64(len(data))/1e6 / elapsed.Seconds())
}

func BenchmarkDecode1(t *testing.B) {
	x, _ := New(2, 2)
	data := make([]byte, int(1e7))
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	vects, _ := x.Split(data)
	x.Encode(vects)

	vects[1] = make([]byte, 0) //drop the 1-th disk
	start := time.Now()
	x.Reconstruct(vects)
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Decode1 time= %v\n", elapsed)
	fmt.Printf("Decode1 Speed= %v MB/s\n", float64(len(data))/1e6 / elapsed.Seconds())

}

func BenchmarkDecode01(t *testing.B) {
	x, _ := New(2, 2)
	data := make([]byte, int(1e7))
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	vects, _ := x.Split(data)
	x.Encode(vects)

	vects[0] = make([]byte, 0) //drop the 0-th and 1-th disk
	vects[1] = make([]byte, 0)
	start := time.Now()
	x.Reconstruct(vects)
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Decode0+1 time= %v\n", elapsed)
	fmt.Printf("Decode0+1 Speed= %v MB/s\n", float64(len(data))/1e6 / elapsed.Seconds())
}

