package hooks

// Hookable defines the interface for models that support hooks
type Hookable interface {
	BeforeSave() error
	AfterFind() error
	BeforeDelete() error
	AfterDelete() error
}

// CallBeforeSave executes the BeforeSave hook if the model implements Hookable
func CallBeforeSave(model interface{}) error {
	if hookable, ok := model.(Hookable); ok {
		return hookable.BeforeSave()
	}
	return nil
}

// CallAfterFind executes the AfterFind hook if the model implements Hookable
func CallAfterFind(model interface{}) error {
	if hookable, ok := model.(Hookable); ok {
		return hookable.AfterFind()
	}
	return nil
}

// CallBeforeDelete executes the BeforeDelete hook if the model implements Hookable
func CallBeforeDelete(model interface{}) error {
	if hookable, ok := model.(Hookable); ok {
		return hookable.BeforeDelete()
	}
	return nil
}

// CallAfterDelete executes the AfterDelete hook if the model implements Hookable
func CallAfterDelete(model interface{}) error {
	if hookable, ok := model.(Hookable); ok {
		return hookable.AfterDelete()
	}
	return nil
}
