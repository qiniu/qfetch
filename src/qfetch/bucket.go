package qfetch

import (
	"fmt"

	"github.com/qiniu/api.v6/auth/digest"
	"github.com/qiniu/api.v6/rs"
	"github.com/qiniu/rpc"
)

type BucketInfo struct {
	Region string `json:"region"`
}

var (
	BUCKET_RS_HOST = "http://rs.qiniu.com"
)

/*
get bucket info
@param mac
@param bucket - bucket name
@return bucketInfo, err
*/
func GetBucketInfo(mac *digest.Mac, bucket string) (bucketInfo BucketInfo, err error) {
	client := rs.New(mac)
	bucketUri := fmt.Sprintf("%s/bucket/%s", BUCKET_RS_HOST, bucket)
	callErr := client.Conn.Call(nil, &bucketInfo, bucketUri)
	if callErr != nil {
		if v, ok := callErr.(*rpc.ErrorInfo); ok {
			err = fmt.Errorf("code: %d, %s, xreqid: %s", v.Code, v.Err, v.Reqid)
		} else {
			err = callErr
		}
	}
	return
}
