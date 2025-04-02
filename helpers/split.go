package main

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	pb "github.com/paschalolo/grpc/helpers/proto"
)

// serializedSize calculates the serialized size of the msg.
// It returns the number of bytes in which the msg was serialized.
func kserializedSize[M protoreflect.ProtoMessage](msg M) int {
	out, err := proto.Marshal(msg)

	if err != nil {
		log.Fatal(err)
	}

	return len(out)
}

func Kmain() {
	s := &pb.Split{Name: "Packt"}
	sz := serializedSize(s)

	fmt.Printf("With Name: %d\n", sz)

	s.Name = ""
	s.ComplexName = &pb.ComplexName{Name: "Packt"}
	sz = kserializedSize(s)

	fmt.Printf("With ComplexName: %d\n", sz)
}
