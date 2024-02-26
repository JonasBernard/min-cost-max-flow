package util

import (
	"fmt"
	"math"
)

func FilterSlice[T any](slice []T, lambda func(T) bool) (ret []T) {
	for _, s := range slice {
		if lambda(s) {
			ret = append(ret, s)
		}
	}
	return
}

func FilterMapBoth[T comparable, S any](m map[T]S, lambda func(T, S) bool) (ret map[T]S) {
	ret = make(map[T]S)
	for t, v := range m {
		if lambda(t, v) {
			ret[t] = v
		}
	}
	return
}

func MapSlice[T any, K any](slice []T, lambda func(*T) K) (ret []K) {
	for _, s := range slice {
		ret = append(ret, lambda(&s))
	}
	return
}

func FlatMapSlice[T any, K any](slice []T, lambda func(*T) []K) (ret []K) {
	for _, s := range slice {
		returnSlice := lambda(&s)
		ret = append(ret, returnSlice...)
	}
	return
}

/*
Takes a map of float64 values and retruns the attained minimum and +Inf if the map is empty.
*/
func MinMapValue[T comparable](values map[T]float64) (min float64) {
	min = math.Inf(1)
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return
}

/*
Takes a map of float64 values and retruns the attained maximum and -Inf if the map is empty.
*/
func MaxMapValue[T comparable](values map[T]float64) (max float64) {
	max = math.Inf(-1)
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return
}

func MapMapValues[T comparable, R, S any](values map[T]R, lambda func(R) S) (ret map[T]S) {
	ret = make(map[T]S)
	for t, v := range values {
		ret[t] = lambda(v)
	}
	return
}

func MinSlice(values []float64) (min float64) {
	min = math.Inf(1)
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return
}

func ArgMin[T comparable](values map[T](float64)) T {
	// fmt.Printf("Finding argmin of %v\n", values)
	var index T
	expectedMin := MinMapValue(values)
	// fmt.Printf("Attained minimum: %v\n", expectedMin)
	for i, v := range values {
		if v == expectedMin {
			index = i
			// fmt.Printf("Found minimal arg %v\n", i)
			break
		}
	}
	return index
}

func Remove[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func RemoveValue[T comparable](slice []T, value T) []T {
	index := Find(slice, value)
	newSlice := Remove(slice, index)
	return newSlice
}

func Find[T comparable](slice []T, value T) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

func PrintMap[S comparable, T any](value map[S]T) {
	for s, t := range value {
		fmt.Printf("%v -> %v\n", s, t)
	}
}
