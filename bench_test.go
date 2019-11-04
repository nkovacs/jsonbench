package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/francoispqt/gojay"
	jsoniter "github.com/json-iterator/go"
	"github.com/nkovacs/jsonx"
	"github.com/wI2L/jettison"
)

var jsoniterStd = jsoniter.ConfigCompatibleWithStandardLibrary

type simplePayload struct {
	St   int    `json:"st"`
	Sid  int    `json:"sid"`
	Tt   string `json:"tt"`
	Gr   int    `json:"gr"`
	UUID string `json:"uuid"`
	IP   string `json:"ip"`
	Ua   string `json:"ua"`
	Tz   int    `json:"tz"`
	V    bool   `json:"v"`
}

func (*simplePayload) NKeys() int    { return 9 }
func (t *simplePayload) IsNil() bool { return t == nil }

func (t *simplePayload) MarshalJSONObject(enc *gojay.Encoder) {
	enc.AddIntKey("st", t.St)
	enc.AddIntKey("sid", t.Sid)
	enc.AddStringKey("tt", t.Tt)
	enc.AddIntKey("gr", t.Gr)
	enc.AddStringKey("uuid", t.UUID)
	enc.AddStringKey("ip", t.IP)
	enc.AddStringKey("ua", t.Ua)
	enc.AddIntKey("tz", t.Tz)
	enc.AddBoolKey("v", t.V)
}

func BenchmarkSimplePayload(b *testing.B) {
	enc, err := jettison.NewEncoder(reflect.TypeOf(simplePayload{}))
	if err != nil {
		b.Fatal(err)
	}
	if err := enc.Compile(); err != nil {
		b.Fatal(err)
	}
	sp := &simplePayload{
		St:   1,
		Sid:  2,
		Tt:   "TestString",
		Gr:   4,
		UUID: "8f9a65eb-4807-4d57-b6e0-bda5d62f1429",
		IP:   "127.0.0.1",
		Ua:   "Mozilla",
		Tz:   8,
		V:    true,
	}
	b.Run("standard", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			bts, err := json.Marshal(sp)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jsonx", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			bts, err := jsonx.Marshal(sp)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jsonx encoder", func(b *testing.B) {
		var buf bytes.Buffer
		enc := jsonx.NewEncoder(&buf)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := enc.Encode(sp)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(buf.Len()))
			buf.Reset()
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			bts, err := jsoniterStd.Marshal(sp)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("gojay", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			bts, err := gojay.MarshalJSONObject(sp)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jettison NoUTF8Coercion NoHTMLEscaping", func(b *testing.B) {
		var buf bytes.Buffer
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// NoUTF8Coercion and NoHTMLEscaping are used to
			// have a fair comparison with Gojay, which does
			// not coerce strings to valid UTF-8 and doesn't
			// escape HTML characters either.
			// None of the string fields of the SimplePayload
			// type contains HTML characters nor contains invalid
			// UTF-8 byte sequences, so this is fine.
			if err := enc.Encode(sp, &buf, jettison.NoUTF8Coercion(), jettison.NoHTMLEscaping()); err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(buf.Len()))
			buf.Reset()
		}
	})
	b.Run("jettison", func(b *testing.B) {
		var buf bytes.Buffer
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if err := enc.Encode(sp, &buf); err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(buf.Len()))
			buf.Reset()
		}
	})
}

func BenchmarkComplexPayload(b *testing.B) {
	type y struct {
		X string `json:"x"`
	}
	type x struct {
		A  y `json:"a"`
		B1 *y
		B2 *y
		C  []string     `json:"c"`
		D  []int        `json:"d"`
		E  []bool       `json:"e"`
		F  []float32    `json:"f,omitempty"`
		G  []*uint      `json:"g"`
		H  [3]string    `json:"h"`
		I  [1]int       `json:"i,omitempty"`
		J  [0]bool      `json:"j"`
		K  []byte       `json:"k"`
		L  []*int       `json:"l"`
		M1 []y          `json:"m1"`
		M2 *[]y         `json:"m2"`
		N  []*y         `json:"n"`
		O1 [3]*int      `json:"o1"`
		O2 *[3]*bool    `json:"o2,omitempty"`
		P  [3]*y        `json:"p"`
		Q  [][]int      `json:"q"`
		R  [2][2]string `json:"r"`
	}
	enc, err := jettison.NewEncoder(reflect.TypeOf(x{}))
	if err != nil {
		b.Fatal(err)
	}
	if err := enc.Compile(); err != nil {
		b.Fatal(err)
	}
	k := make([]byte, 32)
	if _, err := rand.Read(k); err != nil {
		b.Fatal(err)
	}
	var (
		l1, l2 = 0, 42
		m1, m2 = y{X: "Loreum"}, y{}
	)
	xx := &x{
		A:  y{X: "Loreum"},
		B1: nil,
		B2: &y{X: "Ipsum"},
		C:  []string{"one", "two", "three"},
		D:  []int{1, 2, 3},
		E:  []bool{},
		H:  [3]string{"alpha", "beta", "gamma"},
		I:  [1]int{42},
		K:  k,
		L:  []*int{&l1, &l2, nil},
		M1: []y{m1, m2},
		N:  []*y{&m1, &m2, nil},
		O1: [3]*int{&l1, &l2, nil},
		P:  [3]*y{&m1, &m2, nil},
		Q:  [][]int{{1, 2}, {3, 4}},
		R:  [2][2]string{{"a", "b"}, {"c", "d"}},
	}
	b.Run("encoding/json", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bts, err := json.Marshal(xx)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jsonx", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bts, err := jsonx.Marshal(xx)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jsonx encoder", func(b *testing.B) {
		var buf bytes.Buffer
		enc := jsonx.NewEncoder(&buf)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := enc.Encode(xx)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(buf.Len()))
			buf.Reset()
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bts, err := jsoniterStd.Marshal(xx)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jettison", func(b *testing.B) {
		var buf bytes.Buffer
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.Encode(xx, &buf); err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(buf.Len()))
			buf.Reset()
		}
	})
}

func BenchmarkInterface(b *testing.B) {
	s := "Loreum"
	var iface interface{} = s
	enc, err := jettison.NewEncoder(reflect.TypeOf(iface))
	if err != nil {
		b.Fatal(err)
	}
	b.Run("encoding/json", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bts, err := json.Marshal(iface)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jsonx", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bts, err := jsonx.Marshal(iface)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jsonx encoder", func(b *testing.B) {
		var buf bytes.Buffer
		enc := jsonx.NewEncoder(&buf)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := enc.Encode(iface)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(buf.Len()))
			buf.Reset()
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bts, err := jsoniterStd.Marshal(iface)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jettison", func(b *testing.B) {
		var buf bytes.Buffer
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.Encode(iface, &buf); err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(buf.Len()))
			buf.Reset()
		}
	})
}

func BenchmarkMap(b *testing.B) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	b.Run("encoding/json", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bts, err := json.Marshal(m)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jsonx", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bts, err := jsonx.Marshal(m)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jsonx encoder", func(b *testing.B) {
		var buf bytes.Buffer
		enc := jsonx.NewEncoder(&buf)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := enc.Encode(m)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(buf.Len()))
			buf.Reset()
		}
	})
	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bts, err := jsoniterStd.Marshal(m)
			if err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(len(bts)))
		}
	})
	b.Run("jettison", func(b *testing.B) {
		enc, err := jettison.NewEncoder(reflect.TypeOf(m))
		if err != nil {
			b.Fatal(err)
		}
		b.Run("sort", benchMap(enc, m))
		b.Run("nosort", benchMap(enc, m, jettison.UnsortedMap()))
	})
}

func benchMap(enc *jettison.Encoder, m map[string]int, opts ...jettison.Option) func(b *testing.B) {
	var buf bytes.Buffer
	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := enc.Encode(&m, &buf, opts...); err != nil {
				b.Fatal(err)
			}
			b.SetBytes(int64(buf.Len()))
			buf.Reset()
		}
	}
}
