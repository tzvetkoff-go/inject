package inject

import (
	"reflect"
	"strings"
	"unicode"
)

// Injector represents the dependency injector.
type Injector struct {
	Providers map[string]func(interface{}) interface{}
}

// New returns a new injector.
func New() *Injector {
	return &Injector{
		Providers: map[string]func(interface{}) interface{}{},
	}
}

// Reset resets the injector.
func (i *Injector) Reset() {
	i.Providers = map[string]func(interface{}) interface{}{}
}

// Provide registers a provider function for a name.
// The function takes the object that is going to receive the provided value as argument.
// If the function returns nil, the provided value is not set.
func (i *Injector) Provide(name string, provider func(interface{}) interface{}) {
	i.Providers[name] = provider
}

// ProvideObject simply wraps an object inside a provider function.
func (i *Injector) ProvideObject(name string, object interface{}) {
	i.Providers[name] = func(_ interface{}) interface{} {
		return object
	}
}

// GetObject returns an object by name.
func (i *Injector) GetObject(name string, typ interface{}) interface{} {
	if provider := i.Providers[name]; provider != nil {
		return provider(typ)
	}

	return nil
}

// Inject injects all injectable fields.
func (i *Injector) Inject(objects ...interface{}) error {
	for _, object := range objects {
		// Check object type ...
		objectType := reflect.TypeOf(object)
		if objectType.Kind() != reflect.Ptr || objectType.Elem().Kind() != reflect.Struct {
			return NewErrorNotPointerToStruct(object)
		}

		structType := objectType.Elem()
		objectValue := reflect.ValueOf(object)
		indirectValue := reflect.Indirect(objectValue)

		// Walk fields ...
		for idx := 0; idx < structType.NumField(); idx++ {
			// See if injectable ...
			field := structType.Field(idx)
			if injectTag, ok := field.Tag.Lookup("inject"); ok {
				// Lookup name ...
				name := injectTag
				if name == "" {
					name = field.Type.String()
				}

				// Get field value ...
				fieldValue := indirectValue.Field(idx)

				// Inject ...
				// Find provider ...
				provider, ok := i.Providers[name]
				if !ok {
					return NewErrorNoProviderForName(name, object)
				}

				// Call provider ...
				valueToInject := provider(object)

				// If provider returned nil, do nothing ...
				if valueToInject == nil {
					continue
				}

				// Try a setter method ...
				methodName := "Set" + strings.Title(field.Name)
				if method := objectValue.MethodByName(methodName); method.IsValid() {
					method.Call([]reflect.Value{
						reflect.ValueOf(valueToInject),
					})
				} else if unicode.IsLower(rune(field.Name[0])) {
					return NewErrorCannotSetPrivateField(field.Name, object)
				} else {
					fieldValue.Set(reflect.ValueOf(valueToInject))
				}
			}
		}
	}

	return nil
}
