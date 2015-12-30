package lnwire

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	//	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestFundingSignCompleteEncodeDecode(t *testing.T) {
	var (
		//For debugging, writes to /dev/shm/
		//Maybe in the future do it if you do "go test -v"
		WRITE_FILE = false
		FILENAME   = "/dev/shm/fundingSignComplete.raw"

		//TxID
		txid = new(wire.ShaHash)
		//Reversed when displayed
		txidBytes, _ = hex.DecodeString("fd95c6e5c9d5bcf9cfc7231b6a438e46c518c724d0b04b75cc8fddf84a254e3a")
		_            = copy(txid[:], txidBytes)

		//Funding TX Sig 1
		tx                  = wire.NewMsgTx()
		emptybytes          = new([]byte)
		sig1privKeyBytes, _ = hex.DecodeString("927f5827d75dd2addeb532c0fa5ac9277565f981dd6d0d037b422be5f60bdbef")
		sig1privKey, _      = btcec.PrivKeyFromBytes(btcec.S256(), sig1privKeyBytes)
		sigStr1, _          = txscript.RawTxInSignature(tx, 0, *emptybytes, txscript.SigHashAll, sig1privKey)
		commitSig1, _       = btcec.ParseSignature(sigStr1, btcec.S256())
		//Funding TX Sig 2
		sig2privKeyBytes, _ = hex.DecodeString("8a4ad188f6f4000495b765cfb6ffa591133a73019c45428ddd28f53bab551847")
		sig2privKey, _      = btcec.PrivKeyFromBytes(btcec.S256(), sig2privKeyBytes)
		sigStr2, _          = txscript.RawTxInSignature(tx, 0, *emptybytes, txscript.SigHashAll, sig2privKey)
		commitSig2, _       = btcec.ParseSignature(sigStr2, btcec.S256())
		fundingTXSigs       = append(*new([]btcec.Signature), *commitSig1, *commitSig2)

		//funding response
		fundingSignComplete = &FundingSignComplete{
			ReservationID: uint64(12345678),
			TxID:          txid,
			FundingTXSigs: &fundingTXSigs,
		}
		serializedString  = "0000000000bc614efd95c6e5c9d5bcf9cfc7231b6a438e46c518c724d0b04b75cc8fddf84a254e3a02473045022100e7946d057c0b4cc4d3ea525ba156b429796858ebc543d75a6c6c2cbca732db6902202fea377c1f9fb98cd103cf5a4fba276a074b378d4227d15f5fa6439f1a6685bb4630440220235ee55fed634080089953048c3e3f7dc3a154fd7ad18f31dc08e05b7864608a02203bdd7d4e4d9a8162d4b511faf161f0bb16c45181187125017cd0c620c53876ca"
		serializedMessage = "0709110b000000e6000000b80000000000bc614efd95c6e5c9d5bcf9cfc7231b6a438e46c518c724d0b04b75cc8fddf84a254e3a02473045022100e7946d057c0b4cc4d3ea525ba156b429796858ebc543d75a6c6c2cbca732db6902202fea377c1f9fb98cd103cf5a4fba276a074b378d4227d15f5fa6439f1a6685bb4630440220235ee55fed634080089953048c3e3f7dc3a154fd7ad18f31dc08e05b7864608a02203bdd7d4e4d9a8162d4b511faf161f0bb16c45181187125017cd0c620c53876ca"
	)
	//Test serialization
	b := new(bytes.Buffer)
	err := fundingSignComplete.Encode(b, 0)
	if err != nil {
		t.Error("Serialization error")
		t.Error(err.Error())
	} else {
		t.Logf("Encoded FundingSignComplete: %x\n", b.Bytes())
		//Check if we serialized correctly
		if serializedString != hex.EncodeToString(b.Bytes()) {
			t.Error("Serialization does not match expected")
		}

		//So I can do: hexdump -C /dev/shm/fundingSignComplete.raw
		if WRITE_FILE {
			err = ioutil.WriteFile(FILENAME, b.Bytes(), 0644)
			if err != nil {
				t.Error("File write error")
				t.Error(err.Error())
			}
		}
	}

	//Test deserialization
	//Make a new buffer just to be clean
	c := new(bytes.Buffer)
	c.Write(b.Bytes())

	newFunding := NewFundingSignComplete()
	err = newFunding.Decode(c, 0)
	if err != nil {
		t.Error("Decoding Error")
		t.Error(err.Error())
	} else {
		if !reflect.DeepEqual(newFunding, fundingSignComplete) {
			t.Error("Decoding does not match!")
		}
		//Show the struct
		t.Log(newFunding.String())
	}

	//Test message using Message interface
	//Serialize/Encode
	b = new(bytes.Buffer)
	_, err = WriteMessage(b, fundingSignComplete, uint32(1), wire.TestNet3)
	t.Logf("%x\n", b.Bytes())
	if hex.EncodeToString(b.Bytes()) != serializedMessage {
		t.Error("Message encoding error")
	}
	//Deserialize/Decode
	c = new(bytes.Buffer)
	c.Write(b.Bytes())
	_, msg, _, err := ReadMessage(c, uint32(1), wire.TestNet3)
	if err != nil {
		t.Errorf(err.Error())
	} else {
		if !reflect.DeepEqual(msg, fundingSignComplete) {
			t.Error("Message decoding does not match!")
		}
		t.Logf(msg.String())
	}
}
