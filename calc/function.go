package calc

import (
	"math"
	"strconv"
	"strings"
)

const (
	degToRad = 0.017453292519943295769236907684886127134428718885417 // N[Pi/180, 50]
	radToDeg = 57.295779513082320876798154814105170332405472466564   // N[180/Pi, 50]
)

type Function func(float64) float64
type ConstFunction func() float64

func sin(x float64) float64 {
	return math.Sin(x * degToRad)
}

func cos(x float64) float64 {
	return math.Cos(x * degToRad)
}

func tan(x float64) float64 {
	return math.Tan(x * degToRad)
}

func asin(x float64) float64 {
	return math.Asin(x) * radToDeg
}

func acos(x float64) float64 {
	return math.Acos(x) * radToDeg
}

func atan(x float64) float64 {
	return math.Atan(x) * radToDeg
}

func degToDecString(value string) string {
	return strconv.FormatFloat(degToDec(value), 'f', -1, 64)
}

func dms(d, m, s float64) float64 {
	return d + (m / 60.0) + (s / 3600.0)
}

func degToDec(value string) float64 {
	var d, m, s float64

	parse := func(v, chars string) (float64, string) {
		if x, err := strconv.ParseFloat(v, 64); err == nil {
			return x, ""
		}

		if i := strings.IndexAny(v, chars); i != -1 {
			x, err := strconv.ParseFloat(v[:i], 64)
			if err != nil {
				return 0, ""
			}
			return x, v[i+1:]
		}
		return 0, v
	}
	d, value = parse(value, "d")
	m, value = parse(value, "'m")
	s, value = parse(value, `"s`)

	return dms(d, m, s)
}
