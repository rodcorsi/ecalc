package calc

import (
	"math"
	"math/big"
	"strings"
)

const (
	degToRad = 0.017453292519943295769236907684886127134428718885417 // N[Pi/180, 50]
	radToDeg = 57.295779513082320876798154814105170332405472466564   // N[180/Pi, 50]
)

type Function func(*big.Float) *big.Float
type ConstFunction func() *big.Float

func sin(x *big.Float) *big.Float {
	xf, _ := x.Float64()
	return big.NewFloat(math.Sin(xf * degToRad))
}

func cos(x *big.Float) *big.Float {
	xf, _ := x.Float64()
	return big.NewFloat(math.Cos(xf * degToRad))
}

func tan(x *big.Float) *big.Float {
	xf, _ := x.Float64()
	return big.NewFloat(math.Tan(xf * degToRad))
}

func asin(x *big.Float) *big.Float {
	xf, _ := x.Float64()
	return big.NewFloat(math.Asin(xf) * radToDeg)
}

func acos(x *big.Float) *big.Float {
	xf, _ := x.Float64()
	return big.NewFloat(math.Acos(xf) * radToDeg)
}

func atan(x *big.Float) *big.Float {
	xf, _ := x.Float64()
	return big.NewFloat(math.Atan(xf) * radToDeg)
}

func degToDecString(value string) string {
	return degToDec(value).Text('f', -1)
}

func dms(d, m, s *big.Float) *big.Float {
	m60 := new(big.Float).Quo(m, big.NewFloat(60.0))
	s3600 := new(big.Float).Quo(s, big.NewFloat(3600.0))
	return new(big.Float).Add(d, new(big.Float).Add(m60, s3600))
}

func degToDec(value string) *big.Float {
	var d, m, s *big.Float

	parse := func(v, chars string) (*big.Float, string) {
		if x, _, err := big.ParseFloat(v, 10, 256, big.ToNearestEven); err == nil {
			return x, ""
		}

		if i := strings.IndexAny(v, chars); i != -1 {
			x, _, err := big.ParseFloat(v[:i], 10, 256, big.ToNearestEven)
			if err != nil {
				return big.NewFloat(0), ""
			}
			return x, v[i+1:]
		}
		return big.NewFloat(0), v
	}
	d, value = parse(value, "d")
	m, value = parse(value, "'m")
	s, value = parse(value, `"s`)

	return dms(d, m, s)
}

func bigPow(x, y *big.Float) *big.Float {
	xf, _ := x.Float64()
	yf, _ := y.Float64()
	return big.NewFloat(math.Pow(xf, yf))
}

func bigLog(x *big.Float) *big.Float {
	xf, _ := x.Float64()
	return big.NewFloat(math.Log(xf))
}

func bigAbs(x *big.Float) *big.Float {
	return x.Abs(x)
}

func bigSqrt(x *big.Float) *big.Float {
	return x.Sqrt(x)
}

func bigCbrt(x *big.Float) *big.Float {
	xf, _ := x.Float64()
	return big.NewFloat(math.Cbrt(xf))
}

func bigCeil(x *big.Float) *big.Float {
	xf, _ := x.Float64()
	return big.NewFloat(math.Ceil(xf))
}

func bigFloor(x *big.Float) *big.Float {
	xf, _ := x.Float64()
	return big.NewFloat(math.Floor(xf))
}
