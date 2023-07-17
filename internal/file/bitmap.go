package file

func NewBitmap(size int) *Bitmap {
	data := make([]uint64, (size+63)/64)
	return &Bitmap{data: data, size: size}
}

func ParseBitmap(p []byte, size int) *Bitmap {
	data := make([]uint64, len(p)/8)
	for i := range data {
		data[i] = uint64(p[i*8])<<56 |
			uint64(p[i*8+1])<<48 |
			uint64(p[i*8+2])<<40 |
			uint64(p[i*8+3])<<32 |
			uint64(p[i*8+4])<<24 |
			uint64(p[i*8+5])<<16 |
			uint64(p[i*8+6])<<8 |
			uint64(p[i*8+7])
	}
	return &Bitmap{
		data: data,
		size: size,
	}
}

type Bitmap struct {
	data []uint64
	size int
}

func (b *Bitmap) Set(index int) {
	if index < 0 || index >= b.size {
		panic("out of range")
	}
	wordIndex := index / 64
	bitIndex := uint(index % 64)
	b.data[wordIndex] |= 1 << bitIndex
}

func (b *Bitmap) Clear(index int) {
	if index < 0 || index >= b.size {
		panic("out of range")
	}
	wordIndex := index / 64
	bitIndex := uint(index % 64)
	b.data[wordIndex] &= ^(1 << bitIndex)
}

func (b *Bitmap) Get(index int) bool {
	if index < 0 || index >= b.size {
		panic("out of range")
	}
	wordIndex := index / 64
	bitIndex := uint(index % 64)
	return (b.data[wordIndex] & (1 << bitIndex)) != 0
}

func (b *Bitmap) Size() int {
	return b.size
}

func (b *Bitmap) Serialize() []byte {
	p := make([]byte, len(b.data)*8)
	for i, word := range b.data {
		p[i*8] = byte(word >> 56)
		p[i*8+1] = byte(word >> 48)
		p[i*8+2] = byte(word >> 40)
		p[i*8+3] = byte(word >> 32)
		p[i*8+4] = byte(word >> 24)
		p[i*8+5] = byte(word >> 16)
		p[i*8+6] = byte(word >> 8)
		p[i*8+7] = byte(word)
	}
	return p
}
