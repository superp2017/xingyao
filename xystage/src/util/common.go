package util

///不重复压入[]string
func AppendUniqueString(data *[]string, ID string) {
	exist := false
	for _, v := range *data {
		if v == ID {
			exist = true
			break
		}
	}
	if !exist {
		*data = append(*data, ID)
	}
}
