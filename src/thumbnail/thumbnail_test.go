package thumbnail

import "testing"

type TestParseSizeCase struct {
	Raw  string
	Clip string
	N1   int
	N2   int
	Err  error
}

func TestParseSize(t *testing.T) {
	var caselist = make([]TestParseSizeCase, 0)
	caselist = append(caselist, TestParseSizeCase{"b320_120", "b", 320, 120, nil})
	caselist = append(caselist, TestParseSizeCase{"f120_100", "f", 120, 100, nil})
	caselist = append(caselist, TestParseSizeCase{"u120_20", "u", 120, 20, nil})
	caselist = append(caselist, TestParseSizeCase{"p320", "p", 320, 0, nil})
	caselist = append(caselist, TestParseSizeCase{"w120", "w", 120, 0, nil})
	caselist = append(caselist, TestParseSizeCase{"h120", "h", 120, 0, nil})
	for _, c := range caselist {
		clip, n1, n2, err := parseSize(c.Raw)
		t.Log(c)
		if clip != c.Clip {
			t.Error("clip not ok", clip)
		}
		if n1 != c.N1 {
			t.Error("n1 not ok", n1)
		}
		if n2 != c.N2 {
			t.Error("n2 not ok", n2)
		}
		if err != c.Err {
			t.Error("err not ok", err)
		}
	}
}
func TestParseQuality(t *testing.T) {

}
