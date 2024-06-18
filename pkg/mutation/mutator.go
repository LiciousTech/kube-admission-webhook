package mutation

type Mutator interface {
	ApplyMutation(interface{}) (interface{}, interface{}, error)
	Patch(interface{}, interface{}) ([]byte, error)
}
