package conversation_msg

import (
	"context"
	"encoding/json"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/tools/errs"
)

const (
	sendTaskQueueSize        = 256
	sendChainMaxWait         = 3 * time.Second
	defaultMediaOrderedBytes = 16 * 1024
	minMediaOrderedBytes     = 4 * 1024
	maxMediaOrderedBytes     = 8 * 1024 * 1024
	maxSendEnqueueRetry      = 100
	sendEnqueueRetryInterval = time.Millisecond * 5
)

type sendTask struct {
	ctx       context.Context
	msg       *sdk_struct.MsgStruct
	exec      func(context.Context) (*sdk_struct.MsgStruct, error)
	enqueueAt time.Time
	ordered   bool
	lane      ccontext.SendOrderLane
	seq       int64
	mediaSize int64
	deadline  time.Time
}

type messageSender struct {
	conversation *Conversation
	queue        chan *sendTask
	wg           sync.WaitGroup

	textSeq  atomic.Int64
	mediaSeq atomic.Int64

	estimator *thresholdEstimator
}

func newMessageSender(conversation *Conversation) *messageSender {
	workers := runtime.NumCPU()
	if workers < 4 {
		workers = 4
	}
	ms := &messageSender{
		conversation: conversation,
		queue:        make(chan *sendTask, sendTaskQueueSize),
		estimator:    newThresholdEstimator(),
	}
	for i := 0; i < workers; i++ {
		ms.wg.Add(1)
		go ms.worker()
	}
	return ms
}

func (m *messageSender) submit(task *sendTask) error {
	task.enqueueAt = time.Now()
	m.decorate(task)
	for i := 0; i < maxSendEnqueueRetry; i++ {
		select {
		case m.queue <- task:
			return nil
		default:
			time.Sleep(sendEnqueueRetryInterval)
		}
	}
	return errs.New("send task queue full").Wrap()
}

func (m *messageSender) decorate(task *sendTask) {
	task.ordered = true
	if isMediaContentType(task.msg.ContentType) {
		task.lane = ccontext.SendOrderLaneMedia
		task.mediaSize = estimateMediaSize(task.msg)
		task.ordered = m.shouldKeepMediaOrdered(task.mediaSize)
	} else {
		task.lane = ccontext.SendOrderLaneText
	}

	if task.ordered {
		if task.lane == ccontext.SendOrderLaneText {
			task.seq = m.textSeq.Add(1)
		} else {
			task.seq = m.mediaSeq.Add(1)
		}
		task.deadline = task.enqueueAt.Add(sendChainMaxWait)
		task.ctx = ccontext.WithSendOrderInfo(task.ctx, &ccontext.SendOrderInfo{
			Lane:     task.lane,
			Ordered:  true,
			Seq:      task.seq,
			Deadline: task.deadline,
		})
		m.conversation.LongConnMgr.RegisterSendOrder(task.lane, task.seq, task.deadline)
	} else {
		task.seq = 0
		task.deadline = time.Time{}
	}
}

func (m *messageSender) shouldKeepMediaOrdered(size int64) bool {
	if size <= 0 {
		return true
	}
	return float64(size) <= m.estimator.Current()
}

func (m *messageSender) worker() {
	defer m.wg.Done()
	for task := range m.queue {
		m.runTask(task)
	}
}

func (m *messageSender) runTask(task *sendTask) {
	msg, err := task.exec(task.ctx)
	if task.lane == ccontext.SendOrderLaneMedia && task.mediaSize > 0 && err == nil && task.ordered {
		m.estimator.Update(task.mediaSize, time.Since(task.enqueueAt))
	}
	if err != nil {
		notifySendError(task.ctx, err)
		return
	}
	notifySendSuccess(task.ctx, msg)
}

func notifySendSuccess(ctx context.Context, msg *sdk_struct.MsgStruct) {
	callback, _ := ctx.Value(ccontext.Callback).(open_im_sdk_callback.SendMsgCallBack)
	if callback == nil {
		return
	}
	data, err := json.Marshal(msg)
	if err != nil {
		callback.OnError(sdkerrs.UnknownCode, err.Error())
		return
	}
	callback.OnSuccess(string(data))
}

func notifySendError(ctx context.Context, err error) {
	callback, _ := ctx.Value(ccontext.Callback).(open_im_sdk_callback.SendMsgCallBack)
	if callback == nil {
		return
	}
	if code, ok := err.(errs.CodeError); ok {
		callback.OnError(int32(code.Code()), code.Msg())
		return
	}
	callback.OnError(sdkerrs.UnknownCode, err.Error())
}

type thresholdEstimator struct {
	value float64
}

func newThresholdEstimator() *thresholdEstimator {
	return &thresholdEstimator{value: defaultMediaOrderedBytes}
}

func (t *thresholdEstimator) Current() float64 {
	if t.value <= 0 {
		return defaultMediaOrderedBytes
	}
	if t.value > maxMediaOrderedBytes {
		return maxMediaOrderedBytes
	}
	if t.value < minMediaOrderedBytes {
		return minMediaOrderedBytes
	}
	return t.value
}

func (t *thresholdEstimator) Update(size int64, elapsed time.Duration) {
	if size <= 0 || elapsed <= 0 {
		return
	}
	bytesPerSec := float64(size) / elapsed.Seconds()
	if bytesPerSec <= 0 {
		return
	}
	target := bytesPerSec * sendChainMaxWait.Seconds()
	if target > maxMediaOrderedBytes {
		target = maxMediaOrderedBytes
	}
	if target < minMediaOrderedBytes {
		target = minMediaOrderedBytes
	}
	if t.value <= 0 {
		t.value = target
		return
	}
	t.value = 0.6*target + 0.4*t.value
}

func isMediaContentType(contentType int32) bool {
	switch contentType {
	case constant.Picture, constant.Sound, constant.Video, constant.File:
		return true
	default:
		return false
	}
}

func estimateMediaSize(msg *sdk_struct.MsgStruct) int64 {
	switch msg.ContentType {
	case constant.Picture:
		if msg.PictureElem != nil && msg.PictureElem.SourcePicture != nil {
			return msg.PictureElem.SourcePicture.Size
		}
	case constant.Sound:
		if msg.SoundElem != nil {
			return msg.SoundElem.DataSize
		}
	case constant.Video:
		if msg.VideoElem != nil {
			if msg.VideoElem.VideoSize > 0 {
				return msg.VideoElem.VideoSize
			}
			return msg.VideoElem.SnapshotSize
		}
	case constant.File:
		if msg.FileElem != nil {
			return msg.FileElem.FileSize
		}
	}
	return 0
}
