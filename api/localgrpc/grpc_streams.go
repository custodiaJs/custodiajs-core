package localgrpc

import (
	"fmt"
	"io"

	"github.com/CustodiaJS/custodiajs-core/localgrpcproto"
)

func (s *HostAPIService) CoreToProcessControl(stream localgrpcproto.LocalhostAPIService_CoreToProcessControlServer) error {
	// Es wird auf eintreffende Nachrichten gewartet
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println(req)
	}
}
