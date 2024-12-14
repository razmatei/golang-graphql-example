// Code generated by ModelTags. DO NOT EDIT.
package dev

import "emperror.dev/errors"


// ErrObject3UnsupportedJSONKey will be thrown when an unsupported JSON key will be found in transform function.
var ErrObject3UnsupportedJSONKey = errors.Sentinel("unsupported json key")

// ErrObject3UnsupportedStructKeyName will be thrown when an unsupported structure key will be found in transform function.
var ErrObject3UnsupportedStructKeyName = errors.Sentinel("unsupported struct key")

/*
 * JSON Key Names
 */
// Object3 CreatedAt JSON Key Name
const Object3CreatedAtJSONKeyName = "createdAt"

// Object3 DeletedAt JSON Key Name
const Object3DeletedAtJSONKeyName = "deletedAt"

// Object3 ID JSON Key Name
const Object3IDJSONKeyName = "id"

// Object3 Name JSON Key Name
const Object3NameJSONKeyName = "name"

// Object3 UpdatedAt JSON Key Name
const Object3UpdatedAtJSONKeyName = "updatedAt"

var Object3JSONKeyNameList = []string{
    Object3CreatedAtJSONKeyName,
    Object3DeletedAtJSONKeyName,
    Object3IDJSONKeyName,
    Object3NameJSONKeyName,
    Object3UpdatedAtJSONKeyName,
}

/*
 * Struct Key Names
 */
// Object3 CreatedAt Struct Key Name
const Object3CreatedAtStructKeyName = "CreatedAt"

// Object3 DeletedAt Struct Key Name
const Object3DeletedAtStructKeyName = "DeletedAt"

// Object3 ID Struct Key Name
const Object3IDStructKeyName = "ID"

// Object3 Name Struct Key Name
const Object3NameStructKeyName = "Name"

// Object3 UpdatedAt Struct Key Name
const Object3UpdatedAtStructKeyName = "UpdatedAt"

var Object3StructKeyNameList = []string{
    Object3CreatedAtStructKeyName,
    Object3DeletedAtStructKeyName,
    Object3IDStructKeyName,
    Object3NameStructKeyName,
    Object3UpdatedAtStructKeyName,
}


// Transform Object3 Struct Key Name To JSON Key
func TransformObject3StructKeyNameToJSONKey(structKey string) (string, error) {
	switch structKey {
	case Object3CreatedAtStructKeyName:
		return Object3CreatedAtJSONKeyName, nil
	case Object3DeletedAtStructKeyName:
		return Object3DeletedAtJSONKeyName, nil
	case Object3IDStructKeyName:
		return Object3IDJSONKeyName, nil
	case Object3NameStructKeyName:
		return Object3NameJSONKeyName, nil
	case Object3UpdatedAtStructKeyName:
		return Object3UpdatedAtJSONKeyName, nil
	default:
		return "", errors.WithStack(ErrObject3UnsupportedStructKeyName)
	}
}

// Transform Object3 JSON Key To Struct Key Name
func TransformObject3JSONKeyToStructKeyName(jsonKey string) (string, error) {
	switch jsonKey {
	case Object3CreatedAtJSONKeyName:
		return Object3CreatedAtStructKeyName, nil
	case Object3DeletedAtJSONKeyName:
		return Object3DeletedAtStructKeyName, nil
	case Object3IDJSONKeyName:
		return Object3IDStructKeyName, nil
	case Object3NameJSONKeyName:
		return Object3NameStructKeyName, nil
	case Object3UpdatedAtJSONKeyName:
		return Object3UpdatedAtStructKeyName, nil
	default:
		return "", errors.WithStack(ErrObject3UnsupportedJSONKey)
	}
}

// Transform Object3 JSON Key map To Struct Key Name map
func TransformObject3JSONKeyMapToStructKeyNameMap(
	input map[string]interface{},
	ignoreUnsupportedError bool,
) (map[string]interface{}, error) {
	// Rebuild
	m := map[string]interface{}{}
	// Loop over input
	for k, v := range input {
		r, err := TransformObject3JSONKeyToStructKeyName(k)
		// Check error
		if err != nil {
			// Check if ignore is enabled and error is matching
			if ignoreUnsupportedError && errors.Is(err, ErrObject3UnsupportedJSONKey) {
				// Continue the loop
				continue
			}

			// Return
			return nil, err
		}
		// Save
		m[r] = v
	}

	return m, nil
}

// Transform Object3 Struct Key Name map To JSON Key map
func TransformObject3StructKeyNameMapToJSONKeyMap(
	input map[string]interface{},
	ignoreUnsupportedError bool,
) (map[string]interface{}, error) {
	// Rebuild
	m := map[string]interface{}{}
	// Loop over input
	for k, v := range input {
		r, err := TransformObject3StructKeyNameToJSONKey(k)
		// Check error
		if err != nil {
			// Check if ignore is enabled and error is matching
			if ignoreUnsupportedError && errors.Is(err, ErrObject3UnsupportedStructKeyName) {
				// Continue the loop
				continue
			}

			// Return
			return nil, err
		}
		// Save
		m[r] = v
	}

	return m, nil
}