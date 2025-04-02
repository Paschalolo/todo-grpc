package main

import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func serializedSize[D constraints.Integer, W protoreflect.ProtoMessage](data D, wrapper W) (uintptr, int) {
	out, err := proto.Marshal(wrapper)
	if err != nil {
		log.Fatal(err)
	}
	return unsafe.Sizeof(data), len(out) - 1
}
func main() {
	var data uint64 = 172_057_594_037_927_936
	i32 := wrapperspb.UInt64Value{
		Value: data,
	}
	d, w := serializedSize(data, &i32)
	fmt.Printf("in memory : %d \npb : %d\n", d, w)
}
