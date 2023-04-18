package remote

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	pb "sanHeRecruitment/remote/servicepb"
)

// ToService 1
type ToService interface {
	ToService(int *pb.Request, out *pb.Response) error
}

type HttpToService struct {
	BaseURL string
}

func (h *HttpToService) ToService(in *pb.Request, out *pb.Response) error {
	url := fmt.Sprintf(
		"http://%v/_sanheToservice/",
		h.BaseURL,
	)

	fmt.Println(url)

	reqBody, err := proto.Marshal(in)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("NewRequest fail, url: %s, reqBody: %v, err: %v", url, reqBody, err)
		return err
	}
	httpReq.Header.Add("Content-Type", "application/octet-stream")

	// DO: HTTP请求
	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("do http fail, url: %s, reqBody: %v, err:%v", url, reqBody, err)
		return err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", httpRsp.Status)
	}

	Readbytes, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	if err = proto.Unmarshal(Readbytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}
