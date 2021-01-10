package sct

type color struct {
	r, g, b float64
}

/* cribbed from redshift, but truncated with 500K steps */
var whitepoints = []color{
	{1.00000000, 0.18172716, 0.00000000}, /* 1000K */
	{1.00000000, 0.42322816, 0.00000000},
	{1.00000000, 0.54360078, 0.08679949},
	{1.00000000, 0.64373109, 0.28819679},
	{1.00000000, 0.71976951, 0.42860152},
	{1.00000000, 0.77987699, 0.54642268},
	{1.00000000, 0.82854786, 0.64816570},
	{1.00000000, 0.86860704, 0.73688797},
	{1.00000000, 0.90198230, 0.81465502},
	{1.00000000, 0.93853986, 0.88130458},
	{1.00000000, 0.97107439, 0.94305985},
	{1.00000000, 1.00000000, 1.00000000}, /* 6500K */
	{0.95160805, 0.96983355, 1.00000000},
	{0.91194747, 0.94470005, 1.00000000},
	{0.87906581, 0.92357340, 1.00000000},
	{0.85139976, 0.90559011, 1.00000000},
	{0.82782969, 0.89011714, 1.00000000},
	{0.80753191, 0.87667891, 1.00000000},
	{0.78988728, 0.86491137, 1.00000000}, /* 10000K */
	{0.77442176, 0.85453121, 1.00000000},
}

// SetColorTemp changes the monitor colors to reflect the specified color temperature.
func SetColorTemp(temp int) {
	if temp < 1000 || temp > 10000 {
		temp = 6500
	}
	temp -= 1000
	ratio := float64((temp-1000)%500) / 500.0
	point := whitepoints[temp/500]
	point1 := whitepoints[(temp/500)+1]
	gammar := point.r*(1-ratio) + point1.r*ratio
	gammag := point.g*(1-ratio) + point1.g*ratio
	gammab := point.b*(1-ratio) + point1.b*ratio
	setColorTemp(gammar, gammag, gammab)
}
