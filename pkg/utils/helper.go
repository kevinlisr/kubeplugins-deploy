package utils

func InArray( comm []string,c string) bool {
	for _,i := range comm{
		switch i {
		case "get":
			return true
		case "scale":
			return true
		case "list":
			return true
		default:
			return false
		}
	}
	return false
}