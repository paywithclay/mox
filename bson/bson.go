package bson

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HasZeroID checks if the model has an empty ObjectID
func HasZeroID(model interface{}) bool {
	v := reflect.ValueOf(model).Elem()
	idField := v.FieldByName("ID")

	if idField.IsValid() && idField.Type() == reflect.TypeOf(primitive.ObjectID{}) {
		return idField.Interface().(primitive.ObjectID).IsZero()
	}
	return true
}

// GetID retrieves the _id field from the model
func GetID(model interface{}) interface{} {
	v := reflect.ValueOf(model).Elem()
	idField := v.FieldByName("ID")

	if idField.IsValid() && idField.Type() == reflect.TypeOf(primitive.ObjectID{}) {
		return idField.Interface().(primitive.ObjectID)
	}
	return nil
}
