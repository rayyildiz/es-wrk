package worker

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strings"
)

type dataGenerator struct {
	dict  *dict
	TType reflect.Type
}

// NewGenerator represents the data generator.
func NewGenerator(ttype reflect.Type) (*dataGenerator, error) {
	d, err := NewDictionary()
	if err != nil {
		return nil, fmt.Errorf("could not create dictionary, %v", err)
	}
	e := dataGenerator{d, ttype}

	return &e, nil
}

func (generator *dataGenerator) randomWords(n int) string {
	arr, err := generator.dict.GenerateRandomWords(n)
	if err != nil {
		log.Printf("[ERROR] error occurred during creating random words, %v", err)
		return ""
	}

	return strings.Join(arr, " ")
}

func (generator *dataGenerator) initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			len := rand.Intn(10)
			sli := reflect.MakeSlice(ft.Type, len, len)
			//generator.initializeStruct(ft.Type.Elem(),sli.Elem())
			f.Set(sli)
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			generator.initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			generator.initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)

		case reflect.String:
			fv := generator.randomWords(10)
			f.SetString(fv)
		case reflect.Int:
			fv := int64(rand.Intn(100000))
			f.SetInt(fv)
		case reflect.Float32:
			fv := float64(rand.Float32())
			f.SetFloat(fv)
		case reflect.Float64:
			fv := rand.Float64()
			f.SetFloat(fv)
		case reflect.Bool:
			fv := rand.Int()%2 == 0
			f.SetBool(fv)
		default:
		}
	}
}

func (generator *dataGenerator) getRandomElem() interface{} {
	v := reflect.New(generator.TType)
	generator.initializeStruct(generator.TType, v.Elem())
	c := v.Interface()

	return c
}

// GetRandomElements generates random data.
func (generator *dataGenerator) GetRandomElements(numberOfElems int) []interface{} {
	var elems []interface{}

	for i := 0; i < numberOfElems; i++ {
		p := generator.getRandomElem()

		if p != nil {
			elems = append(elems, p)
		}
	}

	return elems
}
