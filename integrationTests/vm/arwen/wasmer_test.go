package arwen

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/integrationTests/vm"
	"github.com/ElrondNetwork/elrond-go/process/factory"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/assert"
)

var ownerAddressBytes = []byte("12345678901234567890123456789012")

func TestAllowNonFloatingPointSC(t *testing.T) {
	wasmvm, scAddress := deploy(t, "floating_point/non_fp.wasm")

	arguments := make([][]byte, 0)
	vmInput := defaultVMInput(arguments)
	vmInput.CallerAddr = ownerAddressBytes

	callInput := makeCallInput(scAddress, "doSomething", vmInput)
	vmOutput, err := wasmvm.RunSmartContractCall(callInput)
	assert.Nil(t, err)

	assert.Equal(t, vmOutput.ReturnCode, vmcommon.Ok)
	fmt.Printf("VM Return Code: %s\n", vmOutput.ReturnCode)
}

func TestDisallowFloatingPointSC(t *testing.T) {
	wasmvm, scAddress := deploy(t, "floating_point/fp.wasm")

	arguments := make([][]byte, 0)
	vmInput := defaultVMInput(arguments)
	vmInput.CallerAddr = ownerAddressBytes

	callInput := makeCallInput(scAddress, "doSomething", vmInput)
	vmOutput, err := wasmvm.RunSmartContractCall(callInput)
	assert.Nil(t, err)

	assert.Equal(t, vmOutput.ReturnCode, vmcommon.ContractInvalid)
	fmt.Printf("VM Return Code: %s\n", vmOutput.ReturnCode)
}

func TestSCAbortExecution_DontAbort(t *testing.T) {
	wasmvm, scAddress := deploy(t, "misc/test_abort.wasm")

	// Run testFunc with argument 0, which will not abort execution, leading to a
	// call to int64finish(100).
	arguments := make([][]byte, 0)
	arguments = append(arguments, []byte{0x00})

	vmInput := defaultVMInput(arguments)
	vmInput.CallerAddr = ownerAddressBytes

	callInput := makeCallInput(scAddress, "testFunc", vmInput)
	vmOutput, err := wasmvm.RunSmartContractCall(callInput)
	assert.Nil(t, err)

	expectedBytes := []byte{100}
	assert.Equal(t, vmOutput.ReturnCode, vmcommon.Ok)
	assertReturnData(t, vmOutput, expectedBytes)
}

func TestSCAbortExecution_Abort(t *testing.T) {
	wasmvm, scAddress := deploy(t, "misc/test_abort.wasm")

	// Run testFunc with argument 1, which will abort execution, leading to a
	// call to int64finish(99).
	arguments := make([][]byte, 0)
	arguments = append(arguments, []byte{0x01})

	vmInput := defaultVMInput(arguments)
	vmInput.CallerAddr = ownerAddressBytes

	callInput := makeCallInput(scAddress, "testFunc", vmInput)
	vmOutput, err := wasmvm.RunSmartContractCall(callInput)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(vmOutput.ReturnData))
	expectedBytes := []byte{98}
	assert.Equal(t, vmOutput.ReturnCode, vmcommon.Ok)
	assertReturnData(t, vmOutput, expectedBytes)
}

func deploy(t *testing.T, wasm_filename string) (vmcommon.VMExecutionHandler, []byte) {
	ownerNonce := uint64(11)
	ownerBalance := big.NewInt(0xfffffffffffffff)
	ownerBalance.Mul(ownerBalance, big.NewInt(0xffffffff))
	round := uint64(444)
	gasPrice := uint64(1)
	gasLimit := uint64(0xffffffffffffffff)

	scCode, err := getBytecode(wasm_filename)
	assert.Nil(t, err)

	scCodeString := hex.EncodeToString(scCode)

	txProc, accnts, blockChainHook := vm.CreatePreparedTxProcessorAndAccountsWithVMs(t, ownerNonce, ownerAddressBytes, ownerBalance)
	scAddressBytes, _ := blockChainHook.NewAddress(ownerAddressBytes, ownerNonce, factory.ArwenVirtualMachine)

	tx := vm.CreateDeployTx(
		ownerAddressBytes,
		ownerNonce,
		big.NewInt(0),
		gasPrice,
		gasLimit,
		scCodeString+"@"+hex.EncodeToString(factory.ArwenVirtualMachine),
	)
	err = txProc.ProcessTransaction(tx, round)
	assert.Nil(t, err)

	vmContainer, blockChainHook := vm.CreateVMAndBlockchainHook(accnts, nil)
	wasmVM, _ := vmContainer.Get(factory.ArwenVirtualMachine)

	return wasmVM, scAddressBytes
}

func assertReturnData(t *testing.T, vmOutput *vmcommon.VMOutput, expectedBytes []byte) {
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode, vmOutput.ReturnCode)
	if len(vmOutput.ReturnData) == 0 {
		assert.True(t, expectedBytes == nil)
		return
	}
	assert.Equal(t, 1, len(vmOutput.ReturnData))
	returnedBytes := vmOutput.ReturnData[0]

	assert.Equal(t, expectedBytes, returnedBytes)
}

func makeCallInput(scAddress []byte, function string, vmInput vmcommon.VMInput) *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		RecipientAddr: scAddress,
		Function:      function,
		VMInput:       vmInput,
	}
}

func defaultVMInput(arguments [][]byte) vmcommon.VMInput {
	return vmcommon.VMInput{
		CallerAddr:  nil,
		CallValue:   big.NewInt(0),
		GasPrice:    uint64(0),
		GasProvided: uint64(0xffffffffffffffff),
		Arguments:   arguments,
	}
}

func padBytesLeft(source []byte, totalLen int, padByte byte) []byte {
	result := make([]byte, 0)
	paddingSize := totalLen - len(source)
	for i := 0; i < paddingSize; i++ {
		result = append(result, padByte)
	}
	for _, b := range source {
		result = append(result, b)
	}
	return result
}

func reverseBytes(source []byte) []byte {
	length := len(source)
	result := make([]byte, 0)
	for i := length - 1; i >= 0; i-- {
		result = append(result, source[i])
	}
	return result
}
