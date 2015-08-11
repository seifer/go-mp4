package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Sync Sample Box (stss - optional)
//
// Contained in : Sample Table box (stbl)
//
// Status: decoded
//
// This lists all sync samples (key frames for video tracks) in the data. If absent, all samples are sync samples.
type StssBox struct {
	Version      byte
	Flags        [3]byte
	header       [8]byte
	SampleNumber []uint32
}

func DecodeStss(r io.Reader) (Box, error) {
	data, err := readAllO(r)

	if err != nil {
		return nil, err
	}

	c := binary.BigEndian.Uint32(data[4:8])
	b := &StssBox{
		Flags:        [3]byte{data[1], data[2], data[3]},
		Version:      data[0],
		SampleNumber: make([]uint32, c),
	}

	for i := 0; i < int(c); i++ {
		b.SampleNumber[i] = binary.BigEndian.Uint32(data[(8 + 4*i):(12 + 4*i)])
	}

	return b, nil
}

func (b *StssBox) Type() string {
	return "stss"
}

func (b *StssBox) Size() int {
	return BoxHeaderSize + 8 + len(b.SampleNumber)*4
}

func (b *StssBox) Dump() {
	fmt.Println("Key frames:")
	for i, n := range b.SampleNumber {
		fmt.Printf(" #%d : sample #%d\n", i, n)
	}
}

func (b *StssBox) Encode(w io.Writer) error {
	binary.BigEndian.PutUint32(b.header[:4], uint32(b.Size()))
	copy(b.header[4:], b.Type())
	_, err := w.Write(b.header[:])
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint32(buf[4:], uint32(len(b.SampleNumber)))
	for i := range b.SampleNumber {
		binary.BigEndian.PutUint32(buf[8+4*i:], b.SampleNumber[i])
	}
	_, err = w.Write(buf)
	return err
}

// Find closest l-frame
func (b *StssBox) GetClosestSample(sample uint32) uint32 {
	sample++

	if len(b.SampleNumber) == 0 {
		return sample
	}

	if sample < b.SampleNumber[0] {
		return b.SampleNumber[0]
	}

	if sample > b.SampleNumber[len(b.SampleNumber)-1] {
		return b.SampleNumber[len(b.SampleNumber)-1]
	}

	for i := 0; i < len(b.SampleNumber); i++ {
		if b.SampleNumber[i] > sample {
			if b.SampleNumber[i]-sample > sample-b.SampleNumber[i-1] {
				return b.SampleNumber[i-1]
			} else {
				return b.SampleNumber[i]
			}
		}
	}

	return sample
}
