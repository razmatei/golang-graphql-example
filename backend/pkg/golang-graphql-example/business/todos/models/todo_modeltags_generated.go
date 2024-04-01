// Code generated by ModelTags. DO NOT EDIT.
package models

import "emperror.dev/errors"

// Todo CreatedAt Gorm Column Name
const TodoCreatedAtGormColumnName = "created_at"

// Todo DeletedAt Gorm Column Name
const TodoDeletedAtGormColumnName = "deleted_at"

// Todo Done Gorm Column Name
const TodoDoneGormColumnName = "done"

// Todo ID Gorm Column Name
const TodoIDGormColumnName = "id"

// Todo Text Gorm Column Name
const TodoTextGormColumnName = "text"

// Todo UpdatedAt Gorm Column Name
const TodoUpdatedAtGormColumnName = "updated_at"

// Todo CreatedAt JSON Key Name
const TodoCreatedAtJSONKeyName = "CreatedAt"

// Todo DeletedAt JSON Key Name
const TodoDeletedAtJSONKeyName = "DeletedAt"

// Todo Done JSON Key Name
const TodoDoneJSONKeyName = "Done"

// Todo ID JSON Key Name
const TodoIDJSONKeyName = "ID"

// Todo Text JSON Key Name
const TodoTextJSONKeyName = "Text"

// Todo UpdatedAt JSON Key Name
const TodoUpdatedAtJSONKeyName = "UpdatedAt"

// Transform Todo Gorm Column To JSON Key
func TransformTodoGormColumnToJSONKey(gormColumn string) (string, error) {
	switch gormColumn {
	case TodoCreatedAtGormColumnName:
		return TodoCreatedAtJSONKeyName, nil
	case TodoDeletedAtGormColumnName:
		return TodoDeletedAtJSONKeyName, nil
	case TodoDoneGormColumnName:
		return TodoDoneJSONKeyName, nil
	case TodoIDGormColumnName:
		return TodoIDJSONKeyName, nil
	case TodoTextGormColumnName:
		return TodoTextJSONKeyName, nil
	case TodoUpdatedAtGormColumnName:
		return TodoUpdatedAtJSONKeyName, nil
	default:
		return "", errors.New("unsupported gorm column")
	}
}


// Transform Todo JSON Key To Gorm Column
func TransformTodoJSONKeyToGormColumn(jsonKey string) (string, error) {
	switch jsonKey {
	case TodoCreatedAtJSONKeyName:
		return TodoCreatedAtGormColumnName, nil
	case TodoDeletedAtJSONKeyName:
		return TodoDeletedAtGormColumnName, nil
	case TodoDoneJSONKeyName:
		return TodoDoneGormColumnName, nil
	case TodoIDJSONKeyName:
		return TodoIDGormColumnName, nil
	case TodoTextJSONKeyName:
		return TodoTextGormColumnName, nil
	case TodoUpdatedAtJSONKeyName:
		return TodoUpdatedAtGormColumnName, nil
	default:
		return "", errors.New("unsupported json key")
	}
}