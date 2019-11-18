package templates

import (
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

type HTLC struct {
	ContractTemplate
}

// MakeHTLC allows a user to recieve the Algo prior to a deadline (in terms of a round) by proving a knowledge
// of a special value or to forfeit the ability to claim, returing it to the payer.
// This contract is usually used to perform cross-chained atomic swaps
//
// More formally -
// Algos can be transferred under only two circumstances:
// 1. To owner if hash_function(arg_0) = hash_value
// 2. To owner if txn.FirstValid > expiry_round
// ...
//
//Parameters
//----------
// - owner : string an address that can receive the asset after the expiry round
// - receiver: string address to receive Algos
// - hashFunction : string the hash function to be used (must be either sha256 or keccak256)
// - hashImage : string the hash image in base64
// - expiryRound : uint64 the round on which the assets can be transferred back to owner
// - maxFee : uint64 the maximum fee that can be paid to the network by the account
func MakeHTLC(owner, receiver, hashFunction, hashImage string, expiryRound, maxFee uint64) (HTLC, error) {
	var referenceProgram string
	if hashFunction == "sha256" {
		referenceProgram = "ASAECAEACSYDIOaalh5vLV96yGYHkmVSvpgjXtMzY8qIkYu5yTipFbb5IH+DsWV/8fxTuS3BgUih1l38LUsfo9Z3KErd0gASbZBpIP68oLsUSlpOp7Q4pGgayA5soQW8tgf8VlMlyVaV9qITMQEiDjEQIxIQMQcyAxIQMQgkEhAxCSgSLQEpEhAxCSoSMQIlDRAREA=="
	} else if hashFunction == "keccak256" {
		referenceProgram = "ASAECAEACSYDIOaalh5vLV96yGYHkmVSvpgjXtMzY8qIkYu5yTipFbb5IH+DsWV/8fxTuS3BgUih1l38LUsfo9Z3KErd0gASbZBpIP68oLsUSlpOp7Q4pGgayA5soQW8tgf8VlMlyVaV9qITMQEiDjEQIxIQMQcyAxIQMQgkEhAxCSgSLQIpEhAxCSoSMQIlDRAREA=="
	} else {
		return HTLC{}, fmt.Errorf("invalid hash function supplied")
	}
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return HTLC{}, err
	}
	ownerAddr, err := types.DecodeAddress(owner)
	if err != nil {
		return HTLC{}, err
	}
	receiverAddr, err := types.DecodeAddress(receiver)
	if err != nil {
		return HTLC{}, err
	}
	//validate hashImage
	_, err = base64.StdEncoding.DecodeString(hashImage)
	if err != nil {
		return HTLC{}, err
	}
	var referenceOffsets = []uint64{ /*fee*/ 3 /*expiryRound*/, 6 /*receiver*/, 10 /*hashImage*/, 43 /*owner*/, 76}
	injectionVector := []interface{}{maxFee, expiryRound, receiverAddr, hashImage, ownerAddr}
	injectedBytes, err := inject(referenceAsBytes, referenceOffsets, injectionVector)
	if err != nil {
		return HTLC{}, err
	}

	injectedProgram := base64.StdEncoding.EncodeToString(injectedBytes)
	address := crypto.AddressFromProgram(injectedBytes)
	htlc := HTLC{}
	htlc.ContractTemplate = ContractTemplate{address: address.String(), program: injectedProgram}
	return htlc, err
}
