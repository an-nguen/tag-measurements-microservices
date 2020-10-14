package utils

import (
	"math"

	"tag-measurements-microservices/pkg/models"
)

type Point struct {
	x float64
	y float64
}

func DouglasPeuckerMeasurement(list []models.Measurement, epsilon float64, t string) []models.Measurement {
	if list == nil {
		return []models.Measurement{}
	}
	if len(list) == 0 {
		return []models.Measurement{}
	}

	dMax := 0.0
	index := 0
	end := len(list) - 1
	for i := 1; i < end; i++ {
		d := perpendicularDistanceMeasurement(list[i], list[0], list[end], t)
		if d > dMax {
			index = i
			dMax = d
		}
	}

	var res []models.Measurement

	if dMax > epsilon {
		recResults1 := DouglasPeuckerMeasurement(list[:index+1], epsilon, t)
		recResults2 := DouglasPeuckerMeasurement(list[index:end], epsilon, t)

		res = append(res, recResults1[:len(recResults1)-1]...)
		res = append(res, recResults2...)
		if len(res) < 2 {
			panic("Problem assembling output")
		}
	} else {
		res = nil
		res = append(res, list[0], list[end])
	}
	return res
}

func perpendicularDistanceMeasurement(pt models.Measurement, lineStart models.Measurement, lineEnd models.Measurement, t string) float64 {
	dx := float64(lineEnd.Date.Unix() - lineStart.Date.Unix())
	dy := 0.0
	if t == "temperature" {
		dy = lineEnd.Temperature - lineStart.Temperature
	} else if t == "humidity" {
		dy = lineEnd.Humidity - lineStart.Humidity
	} else if t == "batteryVolt" {
		dy = lineEnd.Voltage - lineStart.Voltage
	} else if t == "signal" {
		dy = lineEnd.Signal - lineStart.Signal
	}

	mag := math.Pow(math.Pow(dx, 2.0)+math.Pow(dy, 2.0), 0.5)
	if mag > 0.0 {
		dx /= mag
		dy /= mag
	}

	pvx := float64(pt.Date.Unix() - lineStart.Date.Unix())
	pvy := pt.Temperature - lineStart.Temperature

	pvdot := dx*pvx + dy*pvy

	ax := pvx - pvdot*dx
	ay := pvy - pvdot*dy

	return math.Pow(math.Pow(ax, 2.0)+math.Pow(ay, 2.0), 0.5)
}

func DouglasPeucker(list []Point, epsilon float64) []Point {
	dmax := 0.0
	index := 0
	end := len(list) - 1
	for i := 1; i < end; i++ {
		d := perpendicularDistance(list[i], list[0], list[end])
		if d > dmax {
			index = i
			dmax = d
		}
	}

	var res []Point

	if dmax > epsilon {
		recResults1 := DouglasPeucker(list[:index+1], epsilon)
		recResults2 := DouglasPeucker(list[index:end], epsilon)

		res = append(res, recResults1[:len(recResults1)-1]...)
		res = append(res, recResults2...)
		if len(res) < 2 {
			panic("Problem assembling output")
		}
	} else {
		res = nil
		res = append(res, list[0], list[end])
	}
	return res
}

func perpendicularDistance(pt Point, lineStart Point, lineEnd Point) float64 {

	dx := lineEnd.x - lineStart.x
	dy := lineEnd.y - lineStart.y

	mag := math.Pow(math.Pow(dx, 2.0)+math.Pow(dy, 2.0), 0.5)
	if mag > 0.0 {
		dx /= mag
		dy /= mag
	}

	pvx := pt.x - lineStart.x
	pvy := pt.y - lineStart.y

	pvdot := dx*pvx + dy*pvy

	ax := pvx - pvdot*dx
	ay := pvy - pvdot*dy

	return math.Pow(math.Pow(ax, 2.0)+math.Pow(ay, 2.0), 0.5)
}
