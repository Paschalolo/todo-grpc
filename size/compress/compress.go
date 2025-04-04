package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"time"

	pb "github.com/paschalolo/grpc/proto/todo/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func serializedSize[M protoreflect.ProtoMessage](msg M) (int, int) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	out, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := gz.Write(out); err != nil {
		log.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}
	return len(out), len(b.Bytes())
}

func main() {
	task := &pb.Task{
		Id: 1,
		Description: `Lorem ipsum Morbi erat ex, lacinia nec efficitur eget, sagittis ut orci. Etiam in dolor placerat, pharetra ligula et, bibendum neque. Vestibulum vitae congue lectus, sed ultricies augue. Nam iaculis elit nec velit luctus, vitae rutrum nunc imperdiet. Nunc vel turpis sit amet lectus pellentesque tincidunt. Proin commodo tincidunt enim, at sodales mi dictum ac. Maecenas molestie, metus quis malesuada dictum, leo erat egestas lacus, sit amet tristique urna magna a diam. Donec ultricies dui sit amet mi ornare egestas. Phasellus ultricies lectus non interdum pellentesque. Cras nisi tellus, feugiat sed enim quis, tristique interdum lacus. Sed vel pharetra arcu, ac fermentum neque. Morbi mollis sollicitudin varius. Ut sit amet vulputate velit.

Mauris semper neque quis lacinia volutpat. Aenean vestibulum diam ex, sit amet posuere dolor luctus non. Ut consectetur felis blandit ipsum convallis, non lobortis justo facilisis. Ut vitae velit pulvinar, pharetra libero semper, dignissim urna. Nullam quam quam, viverra eget feugiat a, interdum et erat. Morbi fringilla, eros et consequat iaculis, ligula nunc hendrerit neque, ac tincidunt massa sem vitae tortor. Nunc volutpat massa at dapibus pulvinar. Etiam risus sem, dignissim vel blandit eget, maximus lacinia purus.

Cras luctus scelerisque tortor id cursus. Curabitur a interdum nunc, et vulputate diam. Quisque vitae ornare ligula. Nulla massa diam, pulvinar at diam vel, pretium lacinia erat. Quisque ut eros accumsan, ultricies ipsum nec, varius lacus. Fusce pretium vitae libero sit amet pharetra. Phasellus vel odio consequat, cursus risus et, malesuada nisl. In et semper magna, id facilisis nisi. Vestibulum pharetra tortor sed ligula pulvinar, malesuada ornare metus vehicula. Ut eu feugiat massa. Aliquam ut volutpat mi. Aliquam quis augue consectetur, mattis massa non, faucibus ante. Vestibulum rutrum pretium tortor et lacinia.

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed lobortis nec mauris ac placerat. Cras pulvinar dolor at orci semper hendrerit. Nam elementum leo vitae quam commodo, blandit ultricies diam malesuada. Suspendisse lacinia euismod quam interdum mollis. Pellentesque a eleifend ante. Aliquam tempus ultricies velit, eget consequat magna volutpat vitae. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Mauris pulvinar vestibulum congue. Aliquam et magna ultrices justo condimentum varius.`,
		DueDate: timestamppb.New(time.Now().Add(5 * 24 * time.Hour)),
	}
	o, c := serializedSize(task)
	fmt.Printf("original : %d \n compressed : %d \n", o, c)
}
