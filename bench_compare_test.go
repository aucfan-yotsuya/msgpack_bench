package bench_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	shamaton "github.com/shamaton/msgpack"
	"github.com/shamaton/msgpack_bench/protocmp"
	"github.com/shamaton/zeroformatter"
	"github.com/ugorji/go/codec"
	vmihailenco "github.com/vmihailenco/msgpack"
)

type BenchChild struct {
	Int    int
	String string
}
type BenchMarkStruct struct {
	Int    int
	Uint   uint
	Float  float32
	Double float64
	Bool   bool
	String string
	Array  []int
	Map    map[string]uint
	Child  BenchChild
}

var bench = BenchMarkStruct{
	Int:    -123,
	Uint:   456,
	Float:  1.234,
	Double: 6.789,
	Bool:   true,
	String: "this is text.",
	Array:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
	Map:    map[string]uint{"this": 1, "is": 2, "map": 3},
	Child:  BenchChild{Int: 123456, String: "this is struct of child"},
}

var protobench = &protocmp.BenchMarkStruct{
	Int:     int32(bench.Int),
	Uint:    uint32(bench.Uint),
	Float:   bench.Float,
	Double:  bench.Double,
	Bool:    bench.Bool,
	String_: bench.String,
	Array:   []int32{1, 2, 3, 4, 5, 6, 7, 8, 9},
	Map:     map[string]uint32{"this": 1, "is": 2, "map": 3},
	Child:   &protocmp.BenchChild{Int: 123456, String_: "this is struct of child"},
}

var (
	arrayMsgpackBench []byte
	mapMsgpackBench   []byte
	zeroFmtpackBench  []byte
	jsonPackBench     []byte
	gobPackBench      []byte
	protoPackBench    []byte
)

// for codec
var (
	mhBench = &codec.MsgpackHandle{}
)

func initCompare() {
	// ugorji
	mhBench.MapType = reflect.TypeOf(bench)

	d, err := shamaton.EncodeStructAsArray(bench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	arrayMsgpackBench = d
	d, err = shamaton.EncodeStructAsMap(bench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	mapMsgpackBench = d

	d, err = zeroformatter.Serialize(bench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	zeroFmtpackBench = d

	d, err = json.Marshal(bench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	jsonPackBench = d

	d, err = proto.Marshal(protobench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	protoPackBench = d

	buf := bytes.NewBuffer(nil)
	err = gob.NewEncoder(buf).Encode(bench)
	if err != nil {
		fmt.Println("init err : ", err)
		os.Exit(1)
	}
	gobPackBench = buf.Bytes()

	// check
	check()
}

func check() {
	var mp, arr, vmp, varr, c BenchMarkStruct
	shamaton.DecodeStructAsArray(arrayMsgpackBench, &arr)
	shamaton.DecodeStructAsMap(mapMsgpackBench, &mp)
	vmihailenco.Unmarshal(arrayMsgpackBench, &varr)
	vmihailenco.Unmarshal(mapMsgpackBench, &vmp)
	codec.NewDecoderBytes(mapMsgpackBench, mhBench).Decode(&c)

	if !reflect.DeepEqual(mp, arr) {
		fmt.Println("not equal")
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, varr) {
		fmt.Println("not equal")
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, vmp) {
		fmt.Println("not equal")
		os.Exit(1)
	}
	if !reflect.DeepEqual(mp, c) {
		fmt.Println("not equal")
		os.Exit(1)
	}
}

func BenchmarkCompareDecodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := shamaton.DecodeStructAsMap(mapMsgpackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkCompareDecodeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := vmihailenco.Unmarshal(mapMsgpackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := shamaton.DecodeStructAsArray(arrayMsgpackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
func BenchmarkCompareDecodeArrayVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := vmihailenco.Unmarshal(arrayMsgpackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeUgorji(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		dec := codec.NewDecoderBytes(mapMsgpackBench, mhBench)
		err := dec.Decode(&r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeZeroformatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := zeroformatter.Deserialize(&r, zeroFmtpackBench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		err := json.Unmarshal(jsonPackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r BenchMarkStruct
		buf := bytes.NewBuffer(gobPackBench)
		err := gob.NewDecoder(buf).Decode(&r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareDecodeProtocolBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r protocmp.BenchMarkStruct
		err := proto.Unmarshal(protoPackBench, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

/////////////////////////////////////////////////////////////////

func BenchmarkCompareEncodeShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.EncodeStructAsMap(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := vmihailenco.Marshal(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeArrayShamaton(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := shamaton.EncodeStructAsArray(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeArrayVmihailenco(b *testing.B) {
	for i := 0; i < b.N; i++ {

		var buf bytes.Buffer
		enc := vmihailenco.NewEncoder(&buf).StructAsArray(true)
		err := enc.Encode(bench)

		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeUgorji(b *testing.B) {
	for i := 0; i < b.N; i++ {

		b := []byte{}
		enc := codec.NewEncoderBytes(&b, mhBench)
		err := enc.Encode(bench)

		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeZeroformatter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := zeroformatter.Serialize(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeGob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(nil)
		err := gob.NewEncoder(buf).Encode(bench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkCompareEncodeProtocolBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(protobench)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
