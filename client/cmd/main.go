package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/paschalolo/grpc/client/grpcutil"
	client "github.com/paschalolo/grpc/client/handler/grpc"
	pb "github.com/paschalolo/grpc/proto/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("usage client [IP_ADDR] ")
	}
	addr := args[0]
	conn := grpcutil.ServiceConnection(addr)
	clientCall := client.NewClientCaller(conn)
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Fatalf("unexpected error %v \n", err)
		}
	}(conn)
	fmt.Println("--------------ADD--------------")
	dueDate := time.Now().Add(5 * time.Second)
	id1 := clientCall.AddTask("this is a task ", dueDate)
	id2 := clientCall.AddTask("shooping at the mall  ", dueDate)
	id3 := clientCall.AddTask("Go skating", dueDate)
	id4 := clientCall.AddTask("attend prom ", dueDate)
	id5 := clientCall.AddTask("Bible study", dueDate)
	id6 := clientCall.AddTask("Prayer  ", dueDate)
	fmt.Printf("added task %d", id1)
	fmt.Printf("added task %d", id2)
	fmt.Printf("added task %d", id3)
	fmt.Printf("added task %d", id4)
	fmt.Printf("added task %d", id5)
	fmt.Printf("added task %d", id6)
	fmt.Println("-------------------------------")
	fmt.Println("--------------LIST STREAM--------------")
	clientCall.PrintTasks()
	fmt.Println("-------------------------------")
	fmt.Println("-------------UPDATE TASKS--------------")
	clientCall.UpdateTasks([]*pb.UpdateTasksRequest{
		{Task: &pb.Task{Id: id1, Description: "a better day"}},
		{Task: &pb.Task{Id: id2, DueDate: timestamppb.New(dueDate.Add(5 * time.Hour))}},
		{Task: &pb.Task{Id: id3, Done: true}},
	}...)
	clientCall.PrintTasks()
	fmt.Println("-------------------------------")

}
