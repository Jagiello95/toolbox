package slice

func NewSliceUtil() SliceUtil {
	return SliceUtil{}
}

type SliceUtil struct{}

func (s *SliceUtil) InSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
