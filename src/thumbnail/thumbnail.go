package thumbnail

import (
	"errors"
	"helper"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/image/draw"

	"golang.org/x/net/context"
)

func Thumbnail(ctx context.Context) (context.Context, string, error) {
	var (
		err error
	)
	log.Println("process thumbnail")
	r, _ := helper.Req(ctx)
	if strings.Index(r.URL.Path, "/timg") != 0 {
		return ctx, "through", nil
	}
	clip, n1, n2, err := parseSize(r.FormValue("size"))
	quality, err := parseQuality(r.FormValue("quality"))
	src := r.FormValue("src")
	if src == "" {
		return ctx, "err", errors.New("src is empty")
	}

	resp, err := http.Get(src)
	if err != nil {
		return ctx, "err", err
	}
	defer resp.Body.Close()

	// log.Println(ioutil.ReadAll(resp.Body))
	// return ctx, "end", nil
	res, _ := helper.Res(ctx)
	res.Header().Set("x-debug", "1")
	res.Header().Set("Content-Type", "image/jpeg")
	err = drawImage(resp.Body, clip, n1, n2, quality, res)
	if err != nil {
		return ctx, "draw error", err
	}

	// _, _, _ = clip, quality, src
	// _, _ = n1, n2

	return ctx, "end", nil
}

func parseQuality(str string) (int, error) {
	t, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(t), nil
}

func parseSize(str string) (clip string, n1, n2 int, err error) {
	if len(str) == 0 {
		return "", 0, 0, errors.New("size is empty")
	}
	clip = string(str[0])
	switch clip {
	case "p":
		fallthrough
	case "w":
		fallthrough
	case "h":
		var num uint64
		num, err = strconv.ParseUint(string(str[1:]), 10, 32)
		if err != nil {
			return "", 0, 0, errors.New("not a number")
		}
		n1 = int(num)
		return
	case "b":
		fallthrough
	case "f":
		fallthrough
	case "u":
		var (
			num1, num2 int64
		)
		sepIndex := strings.Index(str, "_")
		if sepIndex == -1 {
			return "", 0, 0, errors.New("not well formated")
		}
		num1, err = strconv.ParseInt(string(str[1:sepIndex]), 10, 32)
		num2, err = strconv.ParseInt(string(str[sepIndex+1:]), 10, 32)
		n1 = int(num1)
		n2 = int(num2)
		return
	default:
		return "", 0, 0, errors.New("unknown clip method")
	}

	return "", 0, 0, nil
}

func drawImage(r io.Reader, clip string, n1, n2, quality int, w io.Writer) error {
	var cw, ch int
	rawImg, err := jpeg.Decode(r)
	if err != nil {
		return err
	}

	var rw, rh int
	rw = rawImg.Bounds().Max.X - rawImg.Bounds().Min.X
	rh = rawImg.Bounds().Max.Y - rawImg.Bounds().Min.Y
	switch clip {
	case "p":
		if rw*rh > n1 {
			x := math.Pow(float64(n1)/float64(rw)/float64(rh), 0.5)
			cw = int(math.Floor(float64(rw) * x))
			ch = int(math.Floor(float64(rh) * x))
		} else {
			cw = rw
			ch = rh
		}
	case "w":
		if rw < n1 {
			cw = rw
			ch = rh
		} else {
			cw = n1
			ch = int(math.Trunc(float64(n1) / float64(rw) * float64(rh)))
		}
	case "h":
		if rh < n1 {
			cw = rw
			ch = rh
		} else {
			ch = n1
			cw = int(math.Trunc(float64(n1) / float64(rh) * float64(rw)))
		}
	case "b":
		if rw < n1 && rh < n2 {
			cw = rw
			ch = rh
		} else {
			if float64(rw)/float64(rh) < float64(n1)/float64(n2) {
				// higher
				ch = n1
				cw = int(math.Trunc(float64(n1) / float64(rh) * float64(rw)))
			} else {
				// fatter
				cw = n1
				ch = int(math.Trunc(float64(n1) / float64(rw) * float64(rh)))
			}
		}
	case "f":
		cw = n1
		ch = n2
	case "u":
		cw = n1
		ch = n2
	default:
	}

	cvs := image.NewRGBA(image.Rect(0, 0, cw, ch))
	// draw.Draw(cvs, cvs.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)
	draw.ApproxBiLinear.Scale(cvs, cvs.Bounds(), rawImg, rawImg.Bounds(), nil)
	// draw.Draw(cvs, cvs.Bounds(), rawImg, image.ZP, draw.Src)

	jpeg.Encode(w, cvs, &jpeg.Options{quality})
	return nil

}
