package samdoc

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrModelMustBePointerToStruct = errors.New("provided model must be a pointer to struct")
	ErrQueryFieldMustBeStruct     = errors.New("query field must be struct")
	ErrInvalidPopIndex            = errors.New("invalid index for popping")
	ErrInvalidField               = errors.New("invalid field")
)

type Structure struct {
	m reflect.Value
}

func NewStructure(model interface{}) (*Structure, error) {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		return nil, ErrModelMustBePointerToStruct
	}
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return nil, ErrModelMustBePointerToStruct
	}

	return &Structure{m: v}, nil
}

func (s *Structure) Get(qry FieldQuery) (string, error) {
	var v = s.m
	var err error
	for len(qry) > 0 {
		v, err = s.getFieldByName(v, qry[0])
		if err != nil {
			return "", err
		}
		qry = qry.Pop(0)
	}
	return s.fieldToString(v)
}

func (s *Structure) getFieldByName(v reflect.Value, field string) (reflect.Value, error) {
	var vo reflect.Value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return vo, ErrQueryFieldMustBeStruct
	}
	f := v.FieldByName(field)
	if !f.IsValid() {
		return f, ErrInvalidField
	}
	return f, nil
}

func (s *Structure) fieldToString(v reflect.Value) (string, error) {
	return fmt.Sprintf("%v", v.Interface()), nil
}

type FieldQuery []string

func (fq FieldQuery) Pop(indx int) FieldQuery {
	if len(fq) == 0 {
		return fq
	}
	if indx < 0 || indx > len(fq)-1 {
		panic(ErrInvalidPopIndex)
	}
	if indx == len(fq)-1 {
		return fq[:indx]
	}
	return append(fq[:indx], fq[indx+1:]...)
}
