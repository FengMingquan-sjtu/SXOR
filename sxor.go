package sxor




import (
	"errors"
)

//copy from https://github.com/klauspost/reedsolomon/blob/master/reedsolomon.go
// Encoder is an interface to encode Reed-Salomon parity sets for your data.
type Encoder interface {
	// Encode parity for a set of data shards.
	// Input is 'shards' containing data shards followed by parity shards.
	// The number of shards must match the number given to New().
	// Each shard is a byte array, and they must all be the same size.
	// The parity shards will always be overwritten and the data shards
	// will remain the same, so it is safe for you to read from the
	// data shards while this is running.
	Encode(shards [][]byte) error

	// Reconstruct will recreate the missing shards if possible.
	//
	// Given a list of shards, some of which contain data, fills in the
	// ones that don't have data.
	//
	// The length of the array must be equal to the total number of shards.
	// You indicate that a shard is missing by setting it to nil or zero-length.
	// If a shard is zero-length but has sufficient capacity, that memory will
	// be used, otherwise a new []byte will be allocated.
	//
	// If there are too few shards to reconstruct the missing
	// ones, ErrTooFewShards will be returned.
	//
	// The reconstructed shard set is complete, but integrity is not verified.
	// Use the Verify function to check if data set is ok.
	Reconstruct(shards [][]byte) error

	// ReconstructData will recreate any missing data shards, if possible.
	//
	// Given a list of shards, some of which contain data, fills in the
	// data shards that don't have data.
	//
	// The length of the array must be equal to Shards.
	// You indicate that a shard is missing by setting it to nil or zero-length.
	// If a shard is zero-length but has sufficient capacity, that memory will
	// be used, otherwise a new []byte will be allocated.
	//
	// If there are too few shards to reconstruct the missing
	// ones, ErrTooFewShards will be returned.
	//
	// As the reconstructed shard set may contain missing parity shards,
	// calling the Verify function is likely to fail.
	ReconstructData(shards [][]byte) error

	// Split a data slice into the number of shards given to the encoder,
	// and create empty parity shards.
	//
	// The data will be split into equally sized shards.
	// If the data size isn't dividable by the number of shards,
	// the last shard will contain extra zeros.
	//
	// There must be at least 1 byte otherwise ErrShortData will be
	// returned.
	//
	// The data will not be copied, except for the last shard, so you
	// should not modify the data of the input slice afterwards.
	Split(data []byte) ([][]byte, error)
}

type SXOR struct {
	DataNum int
	ParityNum int
}

// New create coder instance with specific data & parity numbers.
func New(dataNum, parityNum int) (x *SXOR, err error) {
	x = &SXOR{DataNum: dataNum, ParityNum: parityNum}
	return
}

// c <-- a ^ b
func xor(vects [][]byte, idx_a, idx_b, idx_c int){
	size := len(vects[idx_a])
	if len(vects[idx_c]) == 0{
		vects[idx_c] = make([]byte, size)
	}
	for j := 0; j<size; j++{
		vects[idx_c][j] =  vects[idx_a][j] ^ vects[idx_b][j]
	}
	return
}

// Encode encodes data for generating parity.
// Write parity vectors into vects[x.DataNum:].
func (x *SXOR) Encode(vects [][]byte) (err error) {
	for i := 0; i < x.ParityNum; i++ {
		xor(vects, i, i+1, i+x.DataNum)
	}
	return
}


func (x *SXOR) reconst(vects [][]byte, dataOnly bool) error {
	
	updated := false
	for {
		updated = false
		for i := 0; i < x.DataNum; i++ {
			if len(vects[i]) == 0{
				if i>0 && len(vects[i-1]) != 0 && len(vects[x.DataNum+i-1]) !=0 {
					xor(vects, i-1, x.DataNum+i-1, i)
					updated = true
				} else if len(vects[i+1]) != 0 && i<x.ParityNum && len(vects[x.DataNum+i]) !=0{
					xor(vects, i+1, x.DataNum+i, i)
					updated = true
				}
			}
			
		}
		if !updated{
			break
		}
	}

	if ! dataOnly{
		for i:=0; i < x.ParityNum; i++ {
			if len(vects[i+x.DataNum]) == 0{
				if len(vects[i])!=0 && len(vects[i+1])!=0 {
					xor(vects, i, i+1, i+x.DataNum)
				}
			}
		}
	}

	for i := 0; i < x.DataNum+x.ParityNum; i++ {
		if len(vects[i]) != 0 {
			return ErrTooFewShards
		}
	}
	return nil

}


// Wrapping function for MINIO.
// Reconstruct only data vectors, without parity.
func (x *SXOR) ReconstructData(vects [][]byte) error {
	return x.reconst(vects, true)
}

// Wrapping function for MINIO.
// Reconstruct both data and parity
func (x *SXOR) Reconstruct(vects [][]byte) error {
	return x.reconst(vects, false)
}



// Copied from line 1171 of https://github.com/klauspost/reedsolomon/blob/master/reedsolomon.go
// replace: r.DataShards --> x.RS.DataNum,  r.Shards --> TolNum(:=x.DataNum+x.ParityNum)
// Split a data slice into the number of shards given to the encoder,
// and create empty parity shards if necessary.
//
// The data will be split into equally sized shards.
// If the data size isn't divisible by the number of shards,
// the last shard will contain extra zeros.
//
// There must be at least 1 byte otherwise ErrShortData will be
// returned.
//
// The data will not be copied, except for the last shard, so you
// should not modify the data of the input slice afterwards.
func (x *SXOR) Split(data []byte) ([][]byte, error) {
	if len(data) == 0 {
		return nil, ErrShortData
	}
	dataLen := len(data)
	// Calculate number of bytes per data shard.
	perShard := (len(data) + x.DataNum - 1) / x.DataNum


	if cap(data) > len(data) {
		data = data[:cap(data)]
	}

	// Only allocate memory if necessary
	var padding []byte
	TolNum := x.DataNum+x.ParityNum
	if len(data) < (TolNum * perShard) {
		// calculate maximum number of full shards in `data` slice
		fullShards := len(data) / perShard
		padding = make([]byte, TolNum * perShard - perShard * fullShards)
		copy(padding, data[perShard*fullShards:])
		data = data[0 : perShard*fullShards]
	} else {
		for i := dataLen; i < dataLen+x.DataNum; i++ {
			data[i] = 0
		}
	}

	// Split into equal-length shards.
	dst := make([][]byte, TolNum)
	i := 0
	for ; i < len(dst) && len(data) >= perShard; i++ {
		dst[i] = data[:perShard:perShard]
		data = data[perShard:]
	}

	for j := 0; i+j < len(dst); j++ {
		dst[i+j] = padding[:perShard:perShard]
		padding = padding[perShard:]
	}

	return dst, nil
}


// Copied from line 134 of https://github.com/klauspost/reedsolomon/blob/master/reedsolomon.go
// ErrInvShardNum will be returned by New, if you attempt to create
// an Encoder with less than one data shard or less than zero parity
// shards.
var ErrInvShardNum = errors.New("cannot create Encoder with less than one data shard or less than zero parity shards")

// ErrMaxShardNum will be returned by New, if you attempt to create an
// Encoder where data and parity shards are bigger than the order of
// GF(2^8).
var ErrMaxShardNum = errors.New("cannot create Encoder with more than 256 data+parity shards")

// ErrShortData will be returned by Split(), if there isn't enough data
// to fill the number of shards.
var ErrShortData = errors.New("not enough data to fill the number of requested shards")

// ErrTooFewShards is returned if too few shards where given to
// Encode/Verify/Reconstruct/Update. It will also be returned from Reconstruct
// if there were too few shards to reconstruct the missing data.
var ErrTooFewShards = errors.New("too few shards given")
