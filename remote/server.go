package remote

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"sanHeRecruitment/models/websocketModel"
	pb "sanHeRecruitment/remote/servicepb"
	"strings"
)

func (p *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		http.Error(w, errors.New("unexpected path").Error(), http.StatusInternalServerError)
		return
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	//读取请求体
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		body, err := proto.Marshal(&pb.Response{Status: 201, Msg: "read body failed"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(body)
		return
	}

	//解析请求
	Outer := pb.Request{}
	if err = proto.Unmarshal(bytes, &Outer); err != nil {
		body, err := proto.Marshal(&pb.Response{Status: 202, Msg: "proto Unmarshal failed"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(body)
		return
	}

	logicErr := logic(&Outer)

	if logicErr != nil {
		body, err := proto.Marshal(&pb.Response{Status: 203, Msg: "server logic error"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(body)
	}

	//答复

	body, err := proto.Marshal(&pb.Response{Status: 200, Msg: "success"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

func logic(Outer *pb.Request) error {
	websocketModel.TSC.ToServiceMiddleContent <- Outer
	return nil
}
