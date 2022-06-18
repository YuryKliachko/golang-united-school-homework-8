package main

const (
	Add      string = "add"
	List     string = "list"
	FindById string = "findById"
	Remove   string = "remove"
)

var knownOperations = []string{Add, List, FindById, Remove}

func IsValidOperation(operaton string) bool {
	valid := false
	for _, knonwnOperation := range knownOperations {
		if operaton == string(knonwnOperation) {
			valid = true
			break
		}
	}
	return valid
}
