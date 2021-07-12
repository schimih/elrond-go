package vmcommon

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFirstReturnData_VMOutputWithNoReturnDataShouldErr(t *testing.T) {
	vmOutput := VMOutput{
		ReturnData: [][]byte{},
	}

	_, err := vmOutput.GetFirstReturnData(AsBigInt)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no return data")
}

func TestGetFirstReturnData_WithBadReturnDataKindShouldErr(t *testing.T) {
	vmOutput := VMOutput{
		ReturnData: [][]byte{[]byte("100")},
	}

	_, err := vmOutput.GetFirstReturnData(42)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "can't interpret")
}

func TestGetFirstReturnData(t *testing.T) {
	value := big.NewInt(100)

	vmOutput := VMOutput{
		ReturnData: [][]byte{value.Bytes()},
	}

	dataAsBigInt, _ := vmOutput.GetFirstReturnData(AsBigInt)
	dataAsBigIntString, _ := vmOutput.GetFirstReturnData(AsBigIntString)
	dataAsString, _ := vmOutput.GetFirstReturnData(AsString)
	dataAsHex, _ := vmOutput.GetFirstReturnData(AsHex)

	assert.Equal(t, value, dataAsBigInt)
	assert.Equal(t, "100", dataAsBigIntString)
	assert.Equal(t, string(value.Bytes()), dataAsString)
	assert.Equal(t, "64", dataAsHex)
}

func TestOutputContext_MergeCompleteAccounts(t *testing.T) {
	t.Parallel()

	transfer1 := OutputTransfer{
		Value:    big.NewInt(0),
		GasLimit: 9999,
		Data:     []byte("data1"),
	}
	left := &OutputAccount{
		Address:         []byte("addr1"),
		Nonce:           1,
		Balance:         big.NewInt(1000),
		BalanceDelta:    big.NewInt(10000),
		StorageUpdates:  nil,
		Code:            []byte("code1"),
		OutputTransfers: []OutputTransfer{transfer1},
	}
	right := &OutputAccount{
		Address:         []byte("addr2"),
		Nonce:           2,
		Balance:         big.NewInt(2000),
		BalanceDelta:    big.NewInt(20000),
		StorageUpdates:  map[string]*StorageUpdate{"key": {Data: []byte("data"), Offset: []byte("offset")}},
		Code:            []byte("code2"),
		OutputTransfers: []OutputTransfer{transfer1, transfer1},
	}

	expected := &OutputAccount{
		Address:         []byte("addr2"),
		Nonce:           2,
		Balance:         big.NewInt(2000),
		BalanceDelta:    big.NewInt(30000),
		StorageUpdates:  map[string]*StorageUpdate{"key": {Data: []byte("data"), Offset: []byte("offset")}},
		Code:            []byte("code2"),
		OutputTransfers: []OutputTransfer{transfer1, transfer1},
	}

	left.MergeOutputAccounts(right)
	require.Equal(t, expected, left)
}

func TestVMOutput_Size(t *testing.T) {
	t.Parallel()

	vmOutput := &VMOutput{
		ReturnData: [][]byte{
			[]byte("return data 1"),
			[]byte("return data 2"),
		},
		ReturnCode:    Ok,
		ReturnMessage: "unable to carry task",
		GasRemaining:  11384573,
		GasRefund:     big.NewInt(288044),
	}

	fmt.Println(getSize(vmOutput))
}

func getSize(v interface{}) int {
	size := int(reflect.TypeOf(v).Size())
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(v)
		for i := 0; i < s.Len(); i++ {
			size += getSize(s.Index(i).Interface())
		}
	case reflect.Map:
		s := reflect.ValueOf(v)
		keys := s.MapKeys()
		size += int(float64(len(keys)) * 10.79) // approximation from https://golang.org/src/runtime/hashmap.go
		for i := range keys {
			size += getSize(keys[i].Interface()) + getSize(s.MapIndex(keys[i]).Interface())
		}
	case reflect.String:
		size += reflect.ValueOf(v).Len()
	case reflect.Struct:
		s := reflect.ValueOf(v)
		for i := 0; i < s.NumField(); i++ {
			if s.Field(i).CanInterface() {
				size += getSize(s.Field(i).Interface())
			}
		}
	}
	return size
}
