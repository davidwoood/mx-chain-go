package notifier_test

import (
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
	"github.com/ElrondNetwork/elrond-go/outport/mock"
	"github.com/ElrondNetwork/elrond-go/outport/notifier"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	"github.com/ElrondNetwork/elrond-go/testscommon/hashingMocks"
	"github.com/ElrondNetwork/elrond-go/testscommon/marshallerMock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockEventNotifierArgs() notifier.ArgsEventNotifier {
	return notifier.ArgsEventNotifier{
		HttpClient:      &mock.HTTPClientStub{},
		Marshaller:      &marshallerMock.MarshalizerMock{},
		Hasher:          &hashingMocks.HasherMock{},
		PubKeyConverter: &testscommon.PubkeyConverterMock{},
	}
}

func TestNewEventNotifier(t *testing.T) {
	t.Parallel()

	t.Run("nil http client", func(t *testing.T) {
		t.Parallel()

		args := createMockEventNotifierArgs()
		args.HttpClient = nil

		en, err := notifier.NewEventNotifier(args)
		require.Nil(t, en)
		require.Equal(t, notifier.ErrNilHTTPClientWrapper, err)
	})

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		args := createMockEventNotifierArgs()
		args.Marshaller = nil

		en, err := notifier.NewEventNotifier(args)
		require.Nil(t, en)
		require.Equal(t, notifier.ErrNilMarshaller, err)
	})

	t.Run("nil hasher", func(t *testing.T) {
		t.Parallel()

		args := createMockEventNotifierArgs()
		args.Hasher = nil

		en, err := notifier.NewEventNotifier(args)
		require.Nil(t, en)
		require.Equal(t, notifier.ErrNilHasher, err)
	})

	t.Run("nil pub key converter", func(t *testing.T) {
		t.Parallel()

		args := createMockEventNotifierArgs()
		args.PubKeyConverter = nil

		en, err := notifier.NewEventNotifier(args)
		require.Nil(t, en)
		require.Equal(t, notifier.ErrNilPubKeyConverter, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		en, err := notifier.NewEventNotifier(createMockEventNotifierArgs())
		require.Nil(t, err)
		require.NotNil(t, en)
	})
}

func TestSaveBlock(t *testing.T) {
	t.Parallel()

	args := createMockEventNotifierArgs()

	wasCalled := false
	args.HttpClient = &mock.HTTPClientStub{
		PostCalled: func(route string, payload interface{}) error {
			wasCalled = true
			return nil
		},
	}

	en, _ := notifier.NewEventNotifier(args)

	saveBlockData := &outport.ArgsSaveBlockData{
		HeaderHash: []byte{},
		TransactionsPool: &outport.Pool{
			Txs: map[string]data.TransactionHandlerWithGasUsedAndFee{
				"txhash1": nil,
			},
			Scrs: map[string]data.TransactionHandlerWithGasUsedAndFee{
				"scrHash1": nil,
			},
			Logs: []*data.LogData{},
		},
	}

	err := en.SaveBlock(saveBlockData)
	require.Nil(t, err)

	require.True(t, wasCalled)
}

func TestRevertIndexedBlock(t *testing.T) {
	t.Parallel()

	args := createMockEventNotifierArgs()

	wasCalled := false
	args.HttpClient = &mock.HTTPClientStub{
		PostCalled: func(route string, payload interface{}) error {
			wasCalled = true
			return nil
		},
	}

	en, _ := notifier.NewEventNotifier(args)

	header := &block.Header{
		Nonce: 1,
		Round: 2,
		Epoch: 3,
	}
	err := en.RevertIndexedBlock(header, &block.Body{})
	require.Nil(t, err)

	require.True(t, wasCalled)
}

func TestFinalizedBlock(t *testing.T) {
	t.Parallel()

	args := createMockEventNotifierArgs()

	wasCalled := false
	args.HttpClient = &mock.HTTPClientStub{
		PostCalled: func(route string, payload interface{}) error {
			wasCalled = true
			return nil
		},
	}

	en, _ := notifier.NewEventNotifier(args)

	hash := []byte("headerHash")
	err := en.FinalizedBlock(hash)
	require.Nil(t, err)

	require.True(t, wasCalled)
}

func TestMockFunctions(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r != nil {
			assert.Fail(t, fmt.Sprintf("should have not panicked: %v", r))
		}
	}()

	en, err := notifier.NewEventNotifier(createMockEventNotifierArgs())
	require.Nil(t, err)
	require.False(t, en.IsInterfaceNil())

	err = en.SaveRoundsInfo(nil)
	require.Nil(t, err)

	err = en.SaveValidatorsRating("", nil)
	require.Nil(t, err)

	err = en.SaveValidatorsPubKeys(nil, 0)
	require.Nil(t, err)

	err = en.SaveAccounts(0, nil, 0)
	require.Nil(t, err)

	err = en.Close()
	require.Nil(t, err)
}
