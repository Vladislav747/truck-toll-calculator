package service

import (
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"math"
)

type CalculatorServicer interface {
	CalculateDistance(data types.OBUData) (float64, error)
}

type CalculatorService struct {
	prevPoint []float64
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{
		prevPoint: make([]float64, 0),
	}
}

func (s *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	fmt.Println("Calculating distance")
	distance := 0.0
	if len(s.prevPoint) > 0 {
		distance = calculateDistance(s.prevPoint[0], s.prevPoint[1], data.Lat, data.Long)
	}
	s.prevPoint = []float64{data.Lat, data.Long}
	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
