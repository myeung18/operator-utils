package creator

import "fmt"

type NullCreatorImpl struct {
}

func (con *NullCreatorImpl) Create(yamlStr string) (bool, error) {
	//put your logic here
	fmt.Println("Unknown Console type")
	return false, nil;
}

func NewNullCreatorImpl() Creator {
	return new(NullCreatorImpl)
}
