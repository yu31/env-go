package loader

import (
	"encoding/json"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var prefix = "LOADER"

type CustomList struct {
	Name string
	Sex  string
}

func (cu *CustomList) Set(value string) error {
	x := strings.Split(value, "/")
	cu.Name = x[0]
	cu.Sex = x[1]
	return nil
}

type CustomSetter struct {
	Value *string
}

func (cu *CustomSetter) Set(value string) error {
	cu.Value = &value
	return nil
}

type Embedded2 struct {
	Em2Int      int     `loader:"Em2Int"`
	Em2IntPtr   *int    `loader:"Em2IntPtr"`
	Em2Str      string  `loader:"Em2Str"`
	Em2StrPtr   *string `loader:"Em2StrPtr"`
	Em2IntSlice []int   `loader:"Em2IntSlice"`
}

type Embedded1 struct {
	Embedded2 Embedded2 `loader:"Embedded2"`
	Em1Int    int       `loader:"Em1Int"`
	Em1IntPtr *int      `loader:"Em1IntPtr"`
	Em1Str    string    `loader:"Em1Str"`
	Em1StrPtr *string   `loader:"Em1StrPtr"`
}

type Specification struct {
	Bool                bool               `loader:"Bool"`
	BoolPtr             *bool              `loader:"BoolPtr"`
	Str                 string             `loader:"Str,default=x1"`
	StrPtr              *string            `loader:"StrPtr,default=y1"`
	Float32             float32            `loader:"Float32"`
	Float64             float64            `loader:"Float64"`
	Float32Ptr          *float32           `loader:"Float32Ptr"`
	Float64Ptr          *float64           `loader:"Float64Ptr"`
	Int                 int                `loader:"Int"`
	Int8                int8               `loader:"Int8"`
	Int16               int16              `loader:"Int16"`
	Int32               int32              `loader:"Int32"`
	Int64               int64              `loader:"Int64"`
	IntPtr              *int               `loader:"IntPtr"`
	Int8Ptr             *int8              `loader:"Int8Ptr"`
	Int16Ptr            *int16             `loader:"Int16Ptr"`
	Int32Ptr            *int32             `loader:"Int32Ptr"`
	Int64Ptr            *int64             `loader:"Int64Ptr"`
	Uint                uint               `loader:"Uint"`
	Uint8               uint8              `loader:"Uint8"`
	Uint16              uint16             `loader:"Uint16"`
	Uint32              uint32             `loader:"Uint32"`
	Uint64              uint64             `loader:"Uint64"`
	UintPtr             *uint              `loader:"UintPtr"`
	Uint8Ptr            *uint8             `loader:"Uint8Ptr"`
	Uint16Ptr           *uint16            `loader:"Uint16Ptr"`
	Uint32Ptr           *uint32            `loader:"Uint32Ptr"`
	Uint64Ptr           *uint64            `loader:"Uint64Ptr"`
	StrArray            [3]string          `loader:"StrArray"`
	StrArrayPtr         *[3]string         `loader:"StrArrayPtr"`
	StrArrayElePtr      [3]*string         `loader:"StrArrayElePtr"`
	StrArrayValElePtr   *[3]*string        `loader:"StrArrayValElePtr"`
	StrSlice            []string           `loader:"StrSlice,default=x1 x2 x3"`
	StrSlicePtr         *[]string          `loader:"StrSlicePtr,default=y1 y2 y3"`
	StrSliceElePtr      []*string          `loader:"StrSliceElePtr"`
	StrSliceValElePtr   *[]*string         `loader:"StrSliceValElePtr"`
	IntArray            [3]int             `loader:"IntArray"`
	IntArrayPtr         *[3]int            `loader:"IntArrayPtr"`
	IntArrayElePtr      [3]*int            `loader:"IntArrayElePtr"`
	IntArrayValElePtr   *[3]*int           `loader:"IntArrayValElePtr"`
	IntSlice            []int              `loader:"IntSlice"`
	IntSlicePtr         *[]int             `loader:"IntSlicePtr"`
	IntSliceElePtr      []*int             `loader:"IntSliceElePtr"`
	IntSliceValElePtr   *[]*int            `loader:"IntSliceValElePtr"`
	UintArray           [3]uint            `loader:"UintArray"`
	UintArrayPtr        *[3]uint           `loader:"UintArrayPtr"`
	UintArrayElePtr     [3]*uint           `loader:"UintArrayElePtr"`
	UintArrayValElePtr  *[3]*uint          `loader:"UintArrayValElePtr"`
	UintSlice           []uint             `loader:"UintSlice"`
	UintSlicePtr        *[]uint            `loader:"UintSlicePtr"`
	UintSliceElePtr     []*uint            `loader:"UintSliceElePtr"`
	UintSliceValElePtr  *[]*uint           `loader:"UintSliceValElePtr"`
	URL                 url.URL            `loader:"URL"`
	URLPtr              *url.URL           `loader:"URLPtr"`
	TimeDuration        time.Duration      `loader:"TimeDuration"`
	TimeDurationPtr     *time.Duration     `loader:"TimeDurationPtr"`
	MapInt              map[int]int        `loader:"MapInt"`
	MapIntPtr           map[int]*int       `loader:"MapIntPtr"`
	MapStr              map[string]string  `loader:"MapStr"`
	MapStrPtr           map[string]*string `loader:"MapStrPtr"`
	CustomSetter        CustomSetter       `loader:"CustomSetter"`
	CustomSetterPtr     *CustomSetter      `loader:"CustomSetterPtr"`
	Embedded1           Embedded1          `loader:"Embedded1"`
	Embedded1Ptr        *Embedded1         `loader:"Embedded1Ptr"`
	TimeTime            time.Time          `loader:"TimeTime"`
	TimeTimePtr         *time.Time         `loader:"TimeTimePtr"`
	Message             string             `loader:"MESSAGE,default=Hello World"` // used to test default value with space
	CustomList          []CustomList       `loader:"CustomList"`
	CustomListPtr       *[]CustomList      `loader:"CustomListPtr"`
	CustomListElePtr    []*CustomList      `loader:"CustomListElePtr"`
	CustomListValElePtr *[]*CustomList     `loader:"CustomListValElePtr"`
}

var SpecEnvs = `
Bool=false
BoolPtr=true
Str=Str
StrPtr=StrPtr
Float32=1.11
Float64=2.22
Float32Ptr=3.33
Float64Ptr=4.44
Int=11
Int8=12
Int16=13
Int32=14
Int64=15
IntPtr=21
Int8Ptr=22
Int16Ptr=23
Int32Ptr=24
Int64Ptr=25
Uint=31
Uint8=32
Uint16=33
Uint32=34
Uint64=35
UintPtr=41
Uint8Ptr=42
Uint16Ptr=43
Uint32Ptr=44
Uint64Ptr=45
StrArray=a1 b1 c1
StrArrayPtr=a2 b2 c2
StrArrayElePtr=a3 b3 c3
StrArrayValElePtr=a4 b4 c4
StrSlice=d1 e1 f1 g1
StrSlicePtr=d2 e2 f2 g2
StrSliceElePtr=d3 e3 f3 g3
StrSliceValElePtr=d4 e4 f4 g4
IntArray=10 20 30
IntArrayPtr=11 21 31
IntArrayElePtr=12 22 32
IntArrayValElePtr=13 23 33
IntSlice=100 200 300 400
IntSlicePtr=101 201 301 401
IntSliceElePtr=102 202 302 402
IntSliceValElePtr=103 203 303 403
UintArray=14 24 34
UintArrayPtr=15 25 35
UintArrayElePtr=16 26 36
UintArrayValElePtr=17 27 37
UintSlice=104 204 304 404
UintSlicePtr=105 205 305 405
UintSliceElePtr=106 206 306 406
UintSliceValElePtr=107 207 307 407
URL=http://127.0.0.1:8080
URLPtr=http://127.0.0.1:9090
TimeDuration=30s
TimeDurationPtr=1m
MapInt=10:100 20:200
MapIntPtr=30:300 40:400
MapStr=a1:b1 a2:b2:b2
MapStrPtr=c1:d1 c1:d2
CustomSetter=CustomSetter
CustomSetterPtr=CustomSetterPtr
Embedded1_Em1Int=10
Embedded1_Em1IntPtr=11
Embedded1_Em1Str=a1
Embedded1_Em1StrPtr=b1
Embedded1_Embedded2_Em2Int=100
Embedded1_Embedded2_Em2IntPtr=110
Embedded1_Embedded2_Em2Str=a100
Embedded1_Embedded2_Em2StrPtr=b100
Embedded1_Embedded2_Em2IntSlice=1000 1001 10002 10003
Embedded1Ptr_Em1Int=20
Embedded1Ptr_Em1IntPtr=21
Embedded1Ptr_Em1Str=a2
Embedded1Ptr_Em1StrPtr=b2
Embedded1Ptr_Embedded2_Em2Int=200
Embedded1Ptr_Embedded2_Em2IntPtr=210
Embedded1Ptr_Embedded2_Em2Str=a200
Embedded1Ptr_Embedded2_Em2StrPtr=b200
Embedded1Ptr_Embedded2_Em2IntSlice=2000 2001 20002 20003
TimeTime=2020-11-18T15:09:42.532851+08:00
TimeTimePrt=2020-11-18T15:09:42.532851+08:00
CustomList=Joe1/man Lisa1/woman
CustomListPtr=Joe2/man Lisa2/woman
CustomListElePtr=Joe3/man Lisa3/woman
CustomListValElePtr=Joe4/man Lisa4/woman
`

func init() {
	SpecEnvs = strings.TrimPrefix(SpecEnvs, "\n")
	SpecEnvs = strings.TrimSuffix(SpecEnvs, "\n")
}

func TestLoader_Load_ByEnv(t *testing.T) {
	os.Clearenv()
	for _, line := range strings.Split(SpecEnvs, "\n") {
		if line == "" {
			continue
		}
		kv := strings.Split(line, "=")
		err := os.Setenv(strings.ToUpper(prefix+"_"+kv[0]), kv[1])
		require.Nil(t, err)
	}

	s := &Specification{}

	l := New(WithPrefix(prefix))
	err := l.Load(s)
	require.Nil(t, err, "%+v", err)

	b, _ := json.MarshalIndent(s, "", "\t")
	_ = b
	//fmt.Println(string(b))

	require.False(t, s.Bool)
	require.True(t, *s.BoolPtr)
	require.Equal(t, "Str", s.Str)
	require.Equal(t, "StrPtr", *s.StrPtr)
	require.Equal(t, float32(1.11), s.Float32)
	require.Equal(t, float64(2.22), s.Float64)
	require.Equal(t, float32(3.33), *s.Float32Ptr)
	require.Equal(t, float64(4.44), *s.Float64Ptr)
	require.Equal(t, int(11), s.Int)
	require.Equal(t, int8(12), s.Int8)
	require.Equal(t, int16(13), s.Int16)
	require.Equal(t, int32(14), s.Int32)
	require.Equal(t, int64(15), s.Int64)
	require.Equal(t, int(21), *s.IntPtr)
	require.Equal(t, int8(22), *s.Int8Ptr)
	require.Equal(t, int16(23), *s.Int16Ptr)
	require.Equal(t, int32(24), *s.Int32Ptr)
	require.Equal(t, int64(25), *s.Int64Ptr)
	require.Equal(t, uint(31), s.Uint)
	require.Equal(t, uint8(32), s.Uint8)
	require.Equal(t, uint16(33), s.Uint16)
	require.Equal(t, uint32(34), s.Uint32)
	require.Equal(t, uint64(35), s.Uint64)
	require.Equal(t, uint(41), *s.UintPtr)
	require.Equal(t, uint8(42), *s.Uint8Ptr)
	require.Equal(t, uint16(43), *s.Uint16Ptr)
	require.Equal(t, uint32(44), *s.Uint32Ptr)
	require.Equal(t, uint64(45), *s.Uint64Ptr)

	require.Equal(t, [3]string{"a1", "b1", "c1"}, s.StrArray)
	require.Equal(t, [3]string{"a2", "b2", "c2"}, *s.StrArrayPtr)
	require.Equal(t, "a3", *s.StrArrayElePtr[0])
	require.Equal(t, "b3", *s.StrArrayElePtr[1])
	require.Equal(t, "c3", *s.StrArrayElePtr[2])
	require.Equal(t, "a4", *((*s.StrArrayValElePtr)[0]))
	require.Equal(t, "b4", *((*s.StrArrayValElePtr)[1]))
	require.Equal(t, "c4", *((*s.StrArrayValElePtr)[2]))

	require.Equal(t, []string{"d1", "e1", "f1", "g1"}, s.StrSlice)
	require.Equal(t, []string{"d2", "e2", "f2", "g2"}, *s.StrSlicePtr)
	require.Equal(t, "d3", *s.StrSliceElePtr[0])
	require.Equal(t, "e3", *s.StrSliceElePtr[1])
	require.Equal(t, "f3", *s.StrSliceElePtr[2])
	require.Equal(t, "g3", *s.StrSliceElePtr[3])
	require.Equal(t, "d4", *((*s.StrSliceValElePtr)[0]))
	require.Equal(t, "e4", *((*s.StrSliceValElePtr)[1]))
	require.Equal(t, "f4", *((*s.StrSliceValElePtr)[2]))
	require.Equal(t, "g4", *((*s.StrSliceValElePtr)[3]))

	require.Equal(t, [3]int{10, 20, 30}, s.IntArray)
	require.Equal(t, [3]int{11, 21, 31}, *s.IntArrayPtr)
	require.Equal(t, 12, *s.IntArrayElePtr[0])
	require.Equal(t, 22, *s.IntArrayElePtr[1])
	require.Equal(t, 32, *s.IntArrayElePtr[2])
	require.Equal(t, 13, *((*s.IntArrayValElePtr)[0]))
	require.Equal(t, 23, *((*s.IntArrayValElePtr)[1]))
	require.Equal(t, 33, *((*s.IntArrayValElePtr)[2]))

	require.Equal(t, []int{100, 200, 300, 400}, s.IntSlice)
	require.Equal(t, []int{101, 201, 301, 401}, *s.IntSlicePtr)
	require.Equal(t, 102, *s.IntSliceElePtr[0])
	require.Equal(t, 202, *s.IntSliceElePtr[1])
	require.Equal(t, 302, *s.IntSliceElePtr[2])
	require.Equal(t, 402, *s.IntSliceElePtr[3])
	require.Equal(t, 103, *((*s.IntSliceValElePtr)[0]))
	require.Equal(t, 203, *((*s.IntSliceValElePtr)[1]))
	require.Equal(t, 303, *((*s.IntSliceValElePtr)[2]))
	require.Equal(t, 403, *((*s.IntSliceValElePtr)[3]))

	require.Equal(t, [3]uint{14, 24, 34}, s.UintArray)
	require.Equal(t, [3]uint{15, 25, 35}, *s.UintArrayPtr)
	require.Equal(t, uint(16), *s.UintArrayElePtr[0])
	require.Equal(t, uint(26), *s.UintArrayElePtr[1])
	require.Equal(t, uint(36), *s.UintArrayElePtr[2])
	require.Equal(t, uint(17), *((*s.UintArrayValElePtr)[0]))
	require.Equal(t, uint(27), *((*s.UintArrayValElePtr)[1]))
	require.Equal(t, uint(37), *((*s.UintArrayValElePtr)[2]))

	require.Equal(t, []uint{104, 204, 304, 404}, s.UintSlice)
	require.Equal(t, []uint{105, 205, 305, 405}, *s.UintSlicePtr)
	require.Equal(t, uint(106), *s.UintSliceElePtr[0])
	require.Equal(t, uint(206), *s.UintSliceElePtr[1])
	require.Equal(t, uint(306), *s.UintSliceElePtr[2])
	require.Equal(t, uint(406), *s.UintSliceElePtr[3])
	require.Equal(t, uint(107), *((*s.UintSliceValElePtr)[0]))
	require.Equal(t, uint(207), *((*s.UintSliceValElePtr)[1]))
	require.Equal(t, uint(307), *((*s.UintSliceValElePtr)[2]))
	require.Equal(t, uint(407), *((*s.UintSliceValElePtr)[3]))

	require.Equal(t, "http://127.0.0.1:8080", s.URL.String())
	require.Equal(t, "http://127.0.0.1:9090", s.URLPtr.String())
	require.Equal(t, time.Second*30, s.TimeDuration)
	require.Equal(t, time.Minute, *s.TimeDurationPtr)

	require.Equal(t, map[int]int{10: 100, 20: 200}, s.MapInt)
	require.Equal(t, 300, *s.MapIntPtr[30])
	require.Equal(t, 400, *s.MapIntPtr[40])
	require.Equal(t, map[string]string{"a1": "b1", "a2": "b2:b2"}, s.MapStr)
	require.Equal(t, "d2", *s.MapStrPtr["c1"])

	require.Equal(t, "CustomSetter", *s.CustomSetter.Value)
	require.Equal(t, "CustomSetterPtr", *s.CustomSetterPtr.Value)

	// Test Embedded
	require.Equal(t, 10, s.Embedded1.Em1Int)
	require.Equal(t, 11, *s.Embedded1.Em1IntPtr)
	require.Equal(t, "a1", s.Embedded1.Em1Str)
	require.Equal(t, "b1", *s.Embedded1.Em1StrPtr)

	require.Equal(t, 100, s.Embedded1.Embedded2.Em2Int)
	require.Equal(t, 110, *s.Embedded1.Embedded2.Em2IntPtr)
	require.Equal(t, "a100", s.Embedded1.Embedded2.Em2Str)
	require.Equal(t, "b100", *s.Embedded1.Embedded2.Em2StrPtr)
	require.Equal(t, []int{1000, 1001, 10002, 10003}, s.Embedded1.Embedded2.Em2IntSlice)

	require.Equal(t, 20, s.Embedded1Ptr.Em1Int)
	require.Equal(t, 21, *s.Embedded1Ptr.Em1IntPtr)
	require.Equal(t, "a2", s.Embedded1Ptr.Em1Str)
	require.Equal(t, "b2", *s.Embedded1Ptr.Em1StrPtr)

	require.Equal(t, 200, s.Embedded1Ptr.Embedded2.Em2Int)
	require.Equal(t, 210, *s.Embedded1Ptr.Embedded2.Em2IntPtr)
	require.Equal(t, "a200", s.Embedded1Ptr.Embedded2.Em2Str)
	require.Equal(t, "b200", *s.Embedded1Ptr.Embedded2.Em2StrPtr)
	require.Equal(t, []int{2000, 2001, 20002, 20003}, s.Embedded1Ptr.Embedded2.Em2IntSlice)

	// Test CustomList.
	require.Equal(t, 2, len(s.CustomList))
	require.Equal(t, "Joe1", s.CustomList[0].Name)
	require.Equal(t, "man", s.CustomList[0].Sex)
	require.Equal(t, "Lisa1", s.CustomList[1].Name)
	require.Equal(t, "woman", s.CustomList[1].Sex)

	require.Equal(t, 2, len(*s.CustomListPtr))
	require.Equal(t, "Joe2", (*s.CustomListPtr)[0].Name)
	require.Equal(t, "man", (*s.CustomListPtr)[0].Sex)
	require.Equal(t, "Lisa2", (*s.CustomListPtr)[1].Name)
	require.Equal(t, "woman", (*s.CustomListPtr)[1].Sex)

	require.Equal(t, 2, len(s.CustomListElePtr))
	require.Equal(t, "Joe3", s.CustomListElePtr[0].Name)
	require.Equal(t, "man", s.CustomListElePtr[0].Sex)
	require.Equal(t, "Lisa3", s.CustomListElePtr[1].Name)
	require.Equal(t, "woman", s.CustomListElePtr[1].Sex)

	require.Equal(t, 2, len(*s.CustomListValElePtr))
	require.Equal(t, "Joe4", (*s.CustomListValElePtr)[0].Name)
	require.Equal(t, "man", (*s.CustomListValElePtr)[0].Sex)
	require.Equal(t, "Lisa4", (*s.CustomListValElePtr)[1].Name)
	require.Equal(t, "woman", (*s.CustomListValElePtr)[1].Sex)
}

func TestLoader_Load_Default(t *testing.T) {
	os.Clearenv()
	s := &Specification{}

	l := New(WithPrefix(prefix))
	err := l.Load(s)
	require.Nil(t, err, "%+v", err)

	require.Equal(t, "x1", s.Str)
	require.Equal(t, "y1", *s.StrPtr)
	require.Equal(t, []string{"x1", "x2", "x3"}, s.StrSlice)
	require.Equal(t, []string{"y1", "y2", "y3"}, *s.StrSlicePtr)
	require.Equal(t, "Hello World", s.Message)
}

func TestLoader_Load_NotPtr(t *testing.T) {
	os.Clearenv()
	s := Specification{}
	l := New(WithPrefix(prefix))
	err := l.Load(s)
	require.NotNil(t, err, "%+v", err)
	require.Equal(t, ErrNotStructPtr, err)
}

func TestLoader_Load_NotStruct(t *testing.T) {
	os.Clearenv()
	var x string
	l := New(WithPrefix(prefix))
	err := l.Load(&x)
	require.NotNil(t, err, "%+v", err)
	require.Equal(t, ErrNotStructPtr, err)
}

func BenchmarkLoader_Load_ByEnv(b *testing.B) {
	os.Clearenv()
	for _, line := range strings.Split(SpecEnvs, "\n") {
		if line == "" {
			continue
		}
		kv := strings.Split(line, "=")
		err := os.Setenv(strings.ToUpper(prefix+"_"+kv[0]), kv[1])
		require.Nil(b, err)
	}
	l := New(WithPrefix(prefix))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var spec Specification
			_ = l.Load(&spec)
		}
	})
}
