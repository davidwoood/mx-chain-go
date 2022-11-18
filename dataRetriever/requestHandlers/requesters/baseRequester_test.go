package requesters

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/dataRetriever/mock"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	dataRetrieverMocks "github.com/ElrondNetwork/elrond-go/testscommon/dataRetriever"
	"github.com/stretchr/testify/assert"
)

func createMockArgBaseRequester() ArgBaseRequester {
	return ArgBaseRequester{
		RequestSender: &dataRetrieverMocks.TopicRequestSenderStub{},
		Marshaller:    &testscommon.MarshalizerStub{},
	}
}

func Test_createBaseRequester(t *testing.T) {
	t.Parallel()

	baseHandler := createBaseRequester(createMockArgBaseRequester())
	assert.False(t, check.IfNilReflect(baseHandler))
}

func Test_checkArgBase(t *testing.T) {
	t.Parallel()

	t.Run("nil request sender should error", func(t *testing.T) {
		t.Parallel()

		err := checkArgBase(ArgBaseRequester{
			RequestSender: nil,
			Marshaller:    &testscommon.MarshalizerStub{},
		})
		assert.Equal(t, err, dataRetriever.ErrNilRequestSender)
	})
	t.Run("nil marshaller should error", func(t *testing.T) {
		t.Parallel()

		err := checkArgBase(ArgBaseRequester{
			RequestSender: &dataRetrieverMocks.TopicRequestSenderStub{},
			Marshaller:    nil,
		})
		assert.Equal(t, err, dataRetriever.ErrNilMarshalizer)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		err := checkArgBase(createMockArgBaseRequester())
		assert.Nil(t, err)
	})
}

func TestBaseRequester_RequestDataFromHash(t *testing.T) {
	t.Parallel()

	providedEpoch := uint32(1234)
	providedHash := []byte("provided hash")
	providedHashes := [][]byte{providedHash}
	wasCalled := false
	requestSender := &dataRetrieverMocks.TopicRequestSenderStub{
		SendOnRequestTopicCalled: func(rd *dataRetriever.RequestData, originalHashes [][]byte) error {
			wasCalled = true
			assert.Equal(t, providedHash, rd.Value)
			assert.Equal(t, providedHashes, originalHashes)
			assert.Equal(t, dataRetriever.HashType, rd.Type)
			assert.Equal(t, providedEpoch, rd.Epoch)
			return nil
		},
	}
	baseHandler := createBaseRequester(ArgBaseRequester{
		RequestSender: requestSender,
		Marshaller:    &testscommon.MarshalizerStub{},
	})
	assert.False(t, check.IfNilReflect(baseHandler))

	assert.Nil(t, baseHandler.RequestDataFromHash(providedHash, providedEpoch))
	assert.True(t, wasCalled)
}

func TestBaseRequester_NumPeersToQuery(t *testing.T) {
	t.Parallel()

	providedIntra := 123
	providedCross := 100
	wasCalled := false
	requestSender := &dataRetrieverMocks.TopicRequestSenderStub{
		SetNumPeersToQueryCalled: func(intra int, cross int) {
			wasCalled = true
			assert.Equal(t, providedIntra, intra)
			assert.Equal(t, providedCross, cross)
		},
		GetNumPeersToQueryCalled: func() (int, int) {
			return providedIntra, providedCross
		},
	}
	baseHandler := createBaseRequester(ArgBaseRequester{
		RequestSender: requestSender,
		Marshaller:    &testscommon.MarshalizerStub{},
	})
	assert.False(t, check.IfNilReflect(baseHandler))

	baseHandler.SetNumPeersToQuery(providedIntra, providedCross)
	assert.True(t, wasCalled)

	intra, cross := baseHandler.NumPeersToQuery()
	assert.Equal(t, providedIntra, intra)
	assert.Equal(t, providedCross, cross)
}

func TestBaseRequester_SetResolverDebugHandler(t *testing.T) {
	t.Parallel()

	providedDebugHandler := &mock.ResolverDebugHandler{}
	requestSender := &dataRetrieverMocks.TopicRequestSenderStub{
		SetResolverDebugHandlerCalled: func(handler dataRetriever.ResolverDebugHandler) error {
			assert.Equal(t, providedDebugHandler, handler)
			return nil
		},
		ResolverDebugHandlerCalled: func() dataRetriever.ResolverDebugHandler {
			return providedDebugHandler
		},
	}
	baseHandler := createBaseRequester(ArgBaseRequester{
		RequestSender: requestSender,
		Marshaller:    &testscommon.MarshalizerStub{},
	})
	assert.False(t, check.IfNilReflect(baseHandler))

	assert.Nil(t, baseHandler.SetResolverDebugHandler(providedDebugHandler))
	assert.Equal(t, providedDebugHandler, requestSender.ResolverDebugHandler())
}
