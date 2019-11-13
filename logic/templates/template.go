package templates

import (
	"encoding/binary"
	"fmt"
	"github.com/algorand/go-algorand-sdk/types"
)

func replace(buf, newBytes []byte, offset, placeholderLength uint64) []byte {
	output := append(buf[:offset], newBytes...)
	return append(output, buf[(offset+placeholderLength):]...)
}

func inject(original []byte, offsets []uint64, values []interface{}) (result []byte, err error) {
	result = original
	if len(offsets) != len(values) {
		err = fmt.Errorf("length of offsets %v does not match length of replacement values %v", len(offsets), len(values))
		return
	}

	for i, value := range values {
		decodedLength := 0

		if valueAsUint, ok := value.(uint64); ok {
			buffer := make([]byte, 1)
			decodedLength = binary.PutUvarint(buffer, valueAsUint)
			result = replace(result, buffer, offsets[i], 1)

			if decodedLength != 0 {
				for j, _ := range offsets {
					offsets[j] = offsets[j] + uint64(decodedLength) - 1
				}
			}
		} else if addressString, ok := value.(string); ok {
			address, err := types.DecodeAddress(addressString)
			if err != nil {
				return
			}
			addressLen := uint64(32)
			addressBytes := make([]byte, addressLen)
			copy(addressBytes, address[:])
			result = replace(result, addressBytes, offsets[i], addressLen)
		}
		if decodedLength != 0 {
			for j, _ := range offsets {
				offsets[j] = offsets[j] + uint64(decodedLength) - 1
			}
		}
	}
	return
}
