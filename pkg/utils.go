package btcplex

import (
	"bytes"
	"github.com/jmhodges/levigo"
)

type KeyValue struct {
	Key   string
	Value string
}

func GetRange(db *levigo.DB, kStart []byte, kEnd []byte) (values []*KeyValue, err error) {
	ro := levigo.NewReadOptions()
	defer ro.Close()

	it := db.NewIterator(ro)
	defer it.Close()

	it.Seek(kStart)
	endBytes := kEnd
	for {
		if it.Valid() {
			if bytes.Compare(it.Key(), endBytes) > 0 {
				return
			}
			values = append(values, &KeyValue{string(it.Key()[:]), string(it.Value()[:])})
			it.Next()
		} else {
			err = it.GetError()
			return
		}
	}

	return
}

func FloatToUint(x float64) uint64 {
	//frep := strconv.FormatFloat(x, 'g', prec, 64)
	//f, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", x), 64)	
	return uint64(int64((x * float64(100000000.0)) + float64(0.5)))
}
