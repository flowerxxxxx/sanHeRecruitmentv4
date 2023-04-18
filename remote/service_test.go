package remote

import (
	"fmt"
	pb "sanHeRecruitment/remote/servicepb"
	"testing"
)

func TestHttpToService_ToService(t *testing.T) {
	var bb = HttpToService{
		BaseURL: "localhost:9999",
	}
	var Outer *pb.Response
	Outer = &pb.Response{}
	in := pb.Request{
		Username: "yanmingyu",
		Content:  "test_content",
	}
	err := bb.ToService(&in, Outer)
	fmt.Println(err)
	fmt.Println(Outer.Msg, Outer.Status)
}
