package sharding

import (
	"fmt"
	"math/big"
	"math/rand"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/hashing/blake2b"
	"github.com/ElrondNetwork/elrond-go/sharding/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewIndexHashedNodesCoordinatorWithRater_NilRaterShouldErr(t *testing.T) {
	nc, _ := NewIndexHashedNodesCoordinator(createArguments())
	ihgs, err := NewIndexHashedNodesCoordinatorWithRater(nc, nil)

	assert.Nil(t, ihgs)
	assert.Equal(t, ErrNilRater, err)
}

func TestNewIndexHashedNodesCoordinatorWithRater_NilNodesCoordinatorShouldErr(t *testing.T) {
	ihgs, err := NewIndexHashedNodesCoordinatorWithRater(nil, &mock.RaterMock{})

	assert.Nil(t, ihgs)
	assert.Equal(t, ErrNilNodesCoordinator, err)
}

func TestNewIndexHashedGroupSelectorWithRater_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	nc, _ := NewIndexHashedNodesCoordinator(createArguments())
	ihgs, err := NewIndexHashedNodesCoordinatorWithRater(nc, &mock.RaterMock{})
	assert.NotNil(t, ihgs)
	assert.Nil(t, err)
}

//------- LoadEligibleList

func TestIndexHashedGroupSelectorWithRater_SetNilEligibleMapShouldErr(t *testing.T) {
	t.Parallel()
	waiting := createDummyNodesMap(2, 1, "waiting")
	nc, _ := NewIndexHashedNodesCoordinator(createArguments())
	ihgs, _ := NewIndexHashedNodesCoordinatorWithRater(nc, &mock.RaterMock{})
	assert.Equal(t, ErrNilInputNodesMap, ihgs.SetNodesPerShards(nil, waiting, 0, true))
}

func TestIndexHashedGroupSelectorWithRater_OkValShouldWork(t *testing.T) {
	t.Parallel()

	eligibleMap := createDummyNodesMap(3, 1, "waiting")
	waitingMap := make(map[uint32][]Validator)
	nodeShuffler := NewXorValidatorsShuffler(3, 3, 0, false)
	epochStartSubscriber := &mock.EpochStartNotifierStub{}
	bootStorer := mock.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ShardConsensusGroupSize: 2,
		MetaConsensusGroupSize:  1,
		Hasher:                  &mock.HasherMock{},
		Shuffler:                nodeShuffler,
		EpochStartSubscriber:    epochStartSubscriber,
		BootStorer:              bootStorer,
		NbShards:                1,
		EligibleNodes:           eligibleMap,
		WaitingNodes:            waitingMap,
		SelfPublicKey:           []byte("test"),
		ConsensusGroupCache:     &mock.NodesCoordinatorCacheMock{},
		ListIndexUpdater:        &mock.ListIndexUpdaterStub{},
	}
	nc, err := NewIndexHashedNodesCoordinator(arguments)
	assert.Nil(t, err)

	ihgs, err := NewIndexHashedNodesCoordinatorWithRater(nc, &mock.RaterMock{})
	assert.Nil(t, err)
	readEligible := ihgs.nodesConfig[0].eligibleMap[0]
	assert.Equal(t, eligibleMap[0], readEligible)
}

//------- functionality tests

func TestIndexHashedGroupSelectorWithRater_ComputeValidatorsGroup1ValidatorShouldCallGetRating(t *testing.T) {
	t.Parallel()

	list := []Validator{
		mock.NewValidatorMock([]byte("pk0"), []byte("addr0")),
	}

	arguments := createArguments()
	arguments.EligibleNodes[0] = list

	raterCalled := false
	rater := &mock.RaterMock{GetRatingCalled: func(string) uint32 {
		raterCalled = true
		return 1
	}}

	nc, _ := NewIndexHashedNodesCoordinator(arguments)
	ihgs, _ := NewIndexHashedNodesCoordinatorWithRater(nc, rater)
	list2, err := ihgs.ComputeConsensusGroup([]byte("randomness"), 0, 0, 0)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(list2))
	assert.Equal(t, true, raterCalled)
}

func TestIndexHashedGroupSelectorWithRater_ComputeExpandedList(t *testing.T) {
	t.Parallel()

	list := []Validator{
		mock.NewValidatorMock([]byte("pk0"), []byte("addr0")),
		mock.NewValidatorMock([]byte("pk1"), []byte("addr1")),
	}

	listMeta := []Validator{
		mock.NewValidatorMock([]byte("pkMeta1"), []byte("addrMeta1")),
		mock.NewValidatorMock([]byte("pkMeta2"), []byte("addrMeta2")),
	}

	eligibleMap := make(map[uint32][]Validator)
	eligibleMap[0] = list
	eligibleMap[core.MetachainShardId] = listMeta
	waitingMap := make(map[uint32][]Validator)
	nodeShuffler := NewXorValidatorsShuffler(2, 2, 0, false)
	epochStartSubscriber := &mock.EpochStartNotifierStub{}
	bootStorer := mock.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ShardConsensusGroupSize: 2,
		MetaConsensusGroupSize:  1,
		Hasher:                  &mock.HasherMock{},
		Shuffler:                nodeShuffler,
		EpochStartSubscriber:    epochStartSubscriber,
		BootStorer:              bootStorer,
		NbShards:                1,
		EligibleNodes:           eligibleMap,
		WaitingNodes:            waitingMap,
		SelfPublicKey:           []byte("key"),
		ConsensusGroupCache:     &mock.NodesCoordinatorCacheMock{},
		ListIndexUpdater:        &mock.ListIndexUpdaterStub{},
	}

	ratingPk0 := uint32(5)
	ratingPk1 := uint32(1)
	rater := &mock.RaterMock{GetRatingCalled: func(pk string) uint32 {
		if pk == "pk0" {
			return ratingPk0
		}
		if pk == "pk1" {
			return ratingPk1
		}
		return 1
	}}

	nc, _ := NewIndexHashedNodesCoordinator(arguments)
	ihgs, _ := NewIndexHashedNodesCoordinatorWithRater(nc, rater)

	eligibleNodes := ihgs.nodesConfig[0].eligibleMap[0]
	expandedList := ihgs.expandEligibleList(eligibleNodes, &ihgs.nodesConfig[0].mutNodesMaps)
	assert.Equal(t, int(ratingPk0+ratingPk1), len(expandedList))

	occurences := make(map[string]uint32, 2)
	occurences["pk0"] = 0
	occurences["pk1"] = 0
	for _, validator := range expandedList {
		occurences[string(validator.PubKey())]++
	}

	assert.Equal(t, ratingPk0, occurences["pk0"])
	assert.Equal(t, ratingPk1, occurences["pk1"])
}

func BenchmarkIndexHashedGroupSelectorWithRater_ComputeValidatorsGroup63of400(b *testing.B) {
	consensusGroupSize := 63
	list := make([]Validator, 0)

	//generate 400 validators
	for i := 0; i < 400; i++ {
		list = append(list, mock.NewValidatorMock([]byte("pk"+strconv.Itoa(i)), []byte("addr"+strconv.Itoa(i))))
	}

	eligibleMap := make(map[uint32][]Validator)
	waitingMap := make(map[uint32][]Validator)
	eligibleMap[0] = list
	nodeShuffler := NewXorValidatorsShuffler(400, 1, 0, false)
	epochStartSubscriber := &mock.EpochStartNotifierStub{}
	bootStorer := mock.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ShardConsensusGroupSize: consensusGroupSize,
		MetaConsensusGroupSize:  1,
		Hasher:                  &mock.HasherMock{},
		Shuffler:                nodeShuffler,
		EpochStartSubscriber:    epochStartSubscriber,
		BootStorer:              bootStorer,
		NbShards:                1,
		EligibleNodes:           eligibleMap,
		WaitingNodes:            waitingMap,
		SelfPublicKey:           []byte("key"),
		ConsensusGroupCache:     &mock.NodesCoordinatorCacheMock{},
		ListIndexUpdater:        &mock.ListIndexUpdaterStub{},
	}
	ihgs, _ := NewIndexHashedNodesCoordinator(arguments)
	ihgsRater, _ := NewIndexHashedNodesCoordinatorWithRater(ihgs, &mock.RaterMock{})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		randomness := strconv.Itoa(0)
		list2, _ := ihgsRater.ComputeConsensusGroup([]byte(randomness), uint64(0), 0, 0)

		assert.Equal(b, consensusGroupSize, len(list2))
	}
}

func TestIndexHashedGroupSelectorWithRater_GetValidatorWithPublicKeyShouldReturnErrNilPubKey(t *testing.T) {
	t.Parallel()

	list := []Validator{
		mock.NewValidatorMock([]byte("pk0"), []byte("addr0")),
	}
	eligibleMap := make(map[uint32][]Validator)
	waitingMap := make(map[uint32][]Validator)
	eligibleMap[0] = list
	eligibleMap[core.MetachainShardId] = list
	nodeShuffler := NewXorValidatorsShuffler(1, 1, 0, false)
	epochStartSubscriber := &mock.EpochStartNotifierStub{}
	bootStorer := mock.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ShardConsensusGroupSize: 1,
		MetaConsensusGroupSize:  1,
		Hasher:                  &mock.HasherMock{},
		Shuffler:                nodeShuffler,
		EpochStartSubscriber:    epochStartSubscriber,
		BootStorer:              bootStorer,
		NbShards:                1,
		EligibleNodes:           eligibleMap,
		WaitingNodes:            waitingMap,
		SelfPublicKey:           []byte("key"),
		ConsensusGroupCache:     &mock.NodesCoordinatorCacheMock{},
		ListIndexUpdater:        &mock.ListIndexUpdaterStub{},
	}
	nc, _ := NewIndexHashedNodesCoordinator(arguments)
	ihgs, _ := NewIndexHashedNodesCoordinatorWithRater(nc, &mock.RaterMock{})

	_, _, err := ihgs.GetValidatorWithPublicKey(nil, 0)
	assert.Equal(t, ErrNilPubKey, err)
}

func TestIndexHashedGroupSelectorWithRater_GetValidatorWithPublicKeyShouldReturnErrValidatorNotFound(t *testing.T) {
	t.Parallel()

	list := []Validator{
		mock.NewValidatorMock([]byte("pk0"), []byte("addr0")),
	}

	eligibleMap := make(map[uint32][]Validator)
	waitingMap := make(map[uint32][]Validator)
	eligibleMap[0] = list
	eligibleMap[core.MetachainShardId] = list
	nodeShuffler := NewXorValidatorsShuffler(1, 1, 0, false)
	epochStartSubscriber := &mock.EpochStartNotifierStub{}
	bootStorer := mock.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ShardConsensusGroupSize: 1,
		MetaConsensusGroupSize:  1,
		Hasher:                  &mock.HasherMock{},
		Shuffler:                nodeShuffler,
		EpochStartSubscriber:    epochStartSubscriber,
		BootStorer:              bootStorer,
		NbShards:                1,
		EligibleNodes:           eligibleMap,
		WaitingNodes:            waitingMap,
		SelfPublicKey:           []byte("key"),
		ConsensusGroupCache:     &mock.NodesCoordinatorCacheMock{},
		ListIndexUpdater:        &mock.ListIndexUpdaterStub{},
	}
	nc, _ := NewIndexHashedNodesCoordinator(arguments)
	ihgs, _ := NewIndexHashedNodesCoordinatorWithRater(nc, &mock.RaterMock{})

	_, _, err := ihgs.GetValidatorWithPublicKey([]byte("pk1"), 0)
	assert.Equal(t, ErrValidatorNotFound, err)
}

func TestIndexHashedGroupSelectorWithRater_GetValidatorWithPublicKeyShouldWork(t *testing.T) {
	t.Parallel()

	listMeta := []Validator{
		mock.NewValidatorMock([]byte("pk0_meta"), []byte("addr0_meta")),
		mock.NewValidatorMock([]byte("pk1_meta"), []byte("addr1_meta")),
		mock.NewValidatorMock([]byte("pk2_meta"), []byte("addr2_meta")),
	}
	listShard0 := []Validator{
		mock.NewValidatorMock([]byte("pk0_shard0"), []byte("addr0_shard0")),
		mock.NewValidatorMock([]byte("pk1_shard0"), []byte("addr1_shard0")),
		mock.NewValidatorMock([]byte("pk2_shard0"), []byte("addr2_shard0")),
	}
	listShard1 := []Validator{
		mock.NewValidatorMock([]byte("pk0_shard1"), []byte("addr0_shard1")),
		mock.NewValidatorMock([]byte("pk1_shard1"), []byte("addr1_shard1")),
		mock.NewValidatorMock([]byte("pk2_shard1"), []byte("addr2_shard1")),
	}

	eligibleMap := make(map[uint32][]Validator)
	waitingMap := make(map[uint32][]Validator)
	nodeShuffler := NewXorValidatorsShuffler(3, 3, 0, false)
	epochStartSubscriber := &mock.EpochStartNotifierStub{}
	bootStorer := mock.NewStorerMock()

	eligibleMap[core.MetachainShardId] = listMeta
	eligibleMap[0] = listShard0
	eligibleMap[1] = listShard1

	arguments := ArgNodesCoordinator{
		ShardConsensusGroupSize: 1,
		MetaConsensusGroupSize:  1,
		Hasher:                  &mock.HasherMock{},
		Shuffler:                nodeShuffler,
		EpochStartSubscriber:    epochStartSubscriber,
		BootStorer:              bootStorer,
		NbShards:                2,
		EligibleNodes:           eligibleMap,
		WaitingNodes:            waitingMap,
		SelfPublicKey:           []byte("key"),
		ConsensusGroupCache:     &mock.NodesCoordinatorCacheMock{},
		ListIndexUpdater:        &mock.ListIndexUpdaterStub{},
	}
	nc, _ := NewIndexHashedNodesCoordinator(arguments)
	ihgs, _ := NewIndexHashedNodesCoordinatorWithRater(nc, &mock.RaterMock{})

	validator, shardId, err := ihgs.GetValidatorWithPublicKey([]byte("pk0_meta"), 0)
	assert.Nil(t, err)
	assert.Equal(t, core.MetachainShardId, shardId)
	assert.Equal(t, []byte("addr0_meta"), validator.Address())

	validator, shardId, err = ihgs.GetValidatorWithPublicKey([]byte("pk1_shard0"), 0)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), shardId)
	assert.Equal(t, []byte("addr1_shard0"), validator.Address())

	validator, shardId, err = ihgs.GetValidatorWithPublicKey([]byte("pk2_shard1"), 0)
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), shardId)
	assert.Equal(t, []byte("addr2_shard1"), validator.Address())
}

func TestIndexHashedGroupSelectorWithRater_GetAllValidatorsPublicKeys(t *testing.T) {
	t.Parallel()

	shardZeroId := uint32(0)
	shardOneId := uint32(1)
	expectedValidatorsPubKeys := map[uint32][][]byte{
		shardZeroId:           {[]byte("pk0_shard0"), []byte("pk1_shard0"), []byte("pk2_shard0")},
		shardOneId:            {[]byte("pk0_shard1"), []byte("pk1_shard1"), []byte("pk2_shard1")},
		core.MetachainShardId: {[]byte("pk0_meta"), []byte("pk1_meta"), []byte("pk2_meta")},
	}

	listMeta := []Validator{
		mock.NewValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][0], []byte("addr0_meta")),
		mock.NewValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][1], []byte("addr1_meta")),
		mock.NewValidatorMock(expectedValidatorsPubKeys[core.MetachainShardId][2], []byte("addr2_meta")),
	}
	listShard0 := []Validator{
		mock.NewValidatorMock(expectedValidatorsPubKeys[shardZeroId][0], []byte("addr0_shard0")),
		mock.NewValidatorMock(expectedValidatorsPubKeys[shardZeroId][1], []byte("addr1_shard0")),
		mock.NewValidatorMock(expectedValidatorsPubKeys[shardZeroId][2], []byte("addr2_shard0")),
	}
	listShard1 := []Validator{
		mock.NewValidatorMock(expectedValidatorsPubKeys[shardOneId][0], []byte("addr0_shard1")),
		mock.NewValidatorMock(expectedValidatorsPubKeys[shardOneId][1], []byte("addr1_shard1")),
		mock.NewValidatorMock(expectedValidatorsPubKeys[shardOneId][2], []byte("addr2_shard1")),
	}

	eligibleMap := make(map[uint32][]Validator)
	waitingMap := make(map[uint32][]Validator)
	nodeShuffler := NewXorValidatorsShuffler(3, 3, 0, false)
	epochStartSubscriber := &mock.EpochStartNotifierStub{}
	bootStorer := mock.NewStorerMock()

	eligibleMap[core.MetachainShardId] = listMeta
	eligibleMap[shardZeroId] = listShard0
	eligibleMap[shardOneId] = listShard1

	arguments := ArgNodesCoordinator{
		ShardConsensusGroupSize: 1,
		MetaConsensusGroupSize:  1,
		Hasher:                  &mock.HasherMock{},
		Shuffler:                nodeShuffler,
		EpochStartSubscriber:    epochStartSubscriber,
		BootStorer:              bootStorer,
		ShardId:                 shardZeroId,
		NbShards:                2,
		EligibleNodes:           eligibleMap,
		WaitingNodes:            waitingMap,
		SelfPublicKey:           []byte("key"),
		ConsensusGroupCache:     &mock.NodesCoordinatorCacheMock{},
		ListIndexUpdater:        &mock.ListIndexUpdaterStub{},
	}

	nc, _ := NewIndexHashedNodesCoordinator(arguments)
	ihgs, err := NewIndexHashedNodesCoordinatorWithRater(nc, &mock.RaterMock{})
	assert.Nil(t, err)

	allValidatorsPublicKeys, err := ihgs.GetAllValidatorsPublicKeys(0)
	assert.Nil(t, err)
	assert.Equal(t, expectedValidatorsPubKeys, allValidatorsPublicKeys)
}

func BenchmarkIndexHashedGroupSelectorWithRater_TestExpandList(b *testing.B) {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)

	fmt.Println(m.HeapAlloc)

	nrNodes := 40000
	ratingSteps := 100
	array := make([]int, nrNodes*ratingSteps)
	for i := 0; i < nrNodes; i++ {
		for j := 0; j < ratingSteps; j++ {
			array[i*ratingSteps+j] = i
		}
	}

	//a := []int{1, 2, 3, 4, 5, 6, 7, 8}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(array), func(i, j int) { array[i], array[j] = array[j], array[i] })
	m2 := runtime.MemStats{}

	runtime.ReadMemStats(&m2)

	fmt.Println(m2.HeapAlloc)
	fmt.Println(fmt.Sprintf("Used %d MB", (m2.HeapAlloc-m.HeapAlloc)/1024/1024))
	//fmt.Print(array[0:100])
}

func BenchmarkIndexHashedGroupSelectorWithRater_TestHashes(b *testing.B) {
	nrElementsInList := int64(4000000)
	nrHashes := 100

	hasher := blake2b.Blake2b{}

	randomBits := ""

	for i := 0; i < nrHashes; i++ {
		randomBits = fmt.Sprintf("%s%d", randomBits, rand.Intn(2))
	}
	//computedListIndex := int64(0)
	for i := 0; i < nrHashes; i++ {
		computedHash := hasher.Compute(randomBits + fmt.Sprintf("%d", i))
		computedLargeIndex := big.NewInt(0)
		computedLargeIndex.SetBytes(computedHash)
		fmt.Println(big.NewInt(0).Mod(computedLargeIndex, big.NewInt(nrElementsInList)).Int64())
	}

	//fmt.Print(array[0:100])
}

func BenchmarkIndexHashedWithRaterGroupSelector_ComputeValidatorsGroup21of400(b *testing.B) {
	consensusGroupSize := 21
	list := make([]Validator, 0)

	//generate 400 validators
	for i := 0; i < 400; i++ {
		list = append(list, mock.NewValidatorMock([]byte("pk"+strconv.Itoa(i)), []byte("addr"+strconv.Itoa(i))))
	}

	eligibleMap := make(map[uint32][]Validator)
	waitingMap := make(map[uint32][]Validator)
	eligibleMap[0] = list
	nodeShuffler := NewXorValidatorsShuffler(400, 1, 0, false)
	epochStartSubscriber := &mock.EpochStartNotifierStub{}
	bootStorer := mock.NewStorerMock()

	arguments := ArgNodesCoordinator{
		ShardConsensusGroupSize: consensusGroupSize,
		MetaConsensusGroupSize:  1,
		Hasher:                  &mock.HasherMock{},
		Shuffler:                nodeShuffler,
		EpochStartSubscriber:    epochStartSubscriber,
		BootStorer:              bootStorer,
		NbShards:                1,
		EligibleNodes:           eligibleMap,
		WaitingNodes:            waitingMap,
		SelfPublicKey:           []byte("key"),
		ConsensusGroupCache:     &mock.NodesCoordinatorCacheMock{},
	}
	ihgs, _ := NewIndexHashedNodesCoordinator(arguments)
	ihgsRater, _ := NewIndexHashedNodesCoordinatorWithRater(ihgs, &mock.RaterMock{})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		randomness := strconv.Itoa(i)
		list2, _ := ihgsRater.ComputeConsensusGroup([]byte(randomness), 0, 0, 0)

		assert.Equal(b, consensusGroupSize, len(list2))
	}
}
