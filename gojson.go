package gojson


func CheckValid(data []byte) error {
	var d decodeState
	err := checkValid(data, &d.scan)
	if err != nil {
		return err
	}
	return nil
}
