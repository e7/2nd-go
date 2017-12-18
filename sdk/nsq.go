package sdk

import "github.com/nsqio/go-nsq"


// 发布消息
func NsqPublishMsg(prdc *nsq.Producer, title string, msg []byte) error {
	err := prdc.Ping()
	if nil != err {
		return err
	}

	prdc.Publish(title, msg)

	return nil
}