package model

type NihaoMetadata struct{
	CarrierRoute          string  `json:"carrier_route"`
	CongressionalDistrict string  `json:"congressional_district"`
	CountyFips            string  `json:"county_fips"`
	CountyName            string  `json:"county_name"`
	Dst                   bool    `json:"dst"`
	ElotSequence          string  `json:"elot_sequence"`
	ElotSort              string  `json:"elot_sort"`
	Latitude              float64 `json:"latitude"`
	Longitude             float64 `json:"longitude"`
	Precision             string  `json:"precision"`
	Rdi                   string  `json:"rdi"`
	RecordType            string  `json:"record_type"`
	TimeZone              string  `json:"time_zone"`
	UtcOffset             int32   `json:"utc_offset"`
	ZipType               string  `json:"zip_type"`
}
func (this *NihaoMetadata) GetCarrierRoute() string {
	if this == nil {
		return ""
	}
	return this.CarrierRoute
}
func (this *NihaoMetadata) SetCarrierRoute(carrierRoute string) {
	if this == nil {
		return
	}
	this.CarrierRoute = carrierRoute
}
func (this *NihaoMetadata) GetCongressionalDistrict() string {
	if this == nil {
		return ""
	}
	return this.CongressionalDistrict
}
func (this *NihaoMetadata) SetCongressionalDistrict(congressionalDistrict string) {
	if this == nil {
		return
	}
	this.CongressionalDistrict = congressionalDistrict
}
func (this *NihaoMetadata) GetCountyFips() string {
	if this == nil {
		return ""
	}
	return this.CountyFips
}
func (this *NihaoMetadata) SetCountyFips(countyFips string) {
	if this == nil {
		return
	}
	this.CountyFips = countyFips
}
func (this *NihaoMetadata) GetCountyName() string {
	if this == nil {
		return ""
	}
	return this.CountyName
}
func (this *NihaoMetadata) SetCountyName(countyName string) {
	if this == nil {
		return
	}
	this.CountyName = countyName
}
func (this *NihaoMetadata) GetDst() bool {
	if this == nil {
		return false
	}
	return this.Dst
}
func (this *NihaoMetadata) SetDst(dst bool) {
	if this == nil {
		return
	}
	this.Dst = dst
}
func (this *NihaoMetadata) GetElotSequence() string {
	if this == nil {
		return ""
	}
	return this.ElotSequence
}
func (this *NihaoMetadata) SetElotSequence(elotSequence string) {
	if this == nil {
		return
	}
	this.ElotSequence = elotSequence
}
func (this *NihaoMetadata) GetElotSort() string {
	if this == nil {
		return ""
	}
	return this.ElotSort
}
func (this *NihaoMetadata) SetElotSort(elotSort string) {
	if this == nil {
		return
	}
	this.ElotSort = elotSort
}
func (this *NihaoMetadata) GetLatitude() float64 {
	if this == nil {
		return 0
	}
	return this.Latitude
}
func (this *NihaoMetadata) SetLatitude(latitude float64) {
	if this == nil {
		return
	}
	this.Latitude = latitude
}
func (this *NihaoMetadata) GetLongitude() float64 {
	if this == nil {
		return 0
	}
	return this.Longitude
}
func (this *NihaoMetadata) SetLongitude(longitude float64) {
	if this == nil {
		return
	}
	this.Longitude = longitude
}
func (this *NihaoMetadata) GetPrecision() string {
	if this == nil {
		return ""
	}
	return this.Precision
}
func (this *NihaoMetadata) SetPrecision(precision string) {
	if this == nil {
		return
	}
	this.Precision = precision
}
func (this *NihaoMetadata) GetRdi() string {
	if this == nil {
		return ""
	}
	return this.Rdi
}
func (this *NihaoMetadata) SetRdi(rdi string) {
	if this == nil {
		return
	}
	this.Rdi = rdi
}
func (this *NihaoMetadata) GetRecordType() string {
	if this == nil {
		return ""
	}
	return this.RecordType
}
func (this *NihaoMetadata) SetRecordType(recordType string) {
	if this == nil {
		return
	}
	this.RecordType = recordType
}
func (this *NihaoMetadata) GetTimeZone() string {
	if this == nil {
		return ""
	}
	return this.TimeZone
}
func (this *NihaoMetadata) SetTimeZone(timeZone string) {
	if this == nil {
		return
	}
	this.TimeZone = timeZone
}
func (this *NihaoMetadata) GetUtcOffset() int32 {
	if this == nil {
		return 0
	}
	return this.UtcOffset
}
func (this *NihaoMetadata) SetUtcOffset(utcOffset int32) {
	if this == nil {
		return
	}
	this.UtcOffset = utcOffset
}
func (this *NihaoMetadata) GetZipType() string {
	if this == nil {
		return ""
	}
	return this.ZipType
}
func (this *NihaoMetadata) SetZipType(zipType string) {
	if this == nil {
		return
	}
	this.ZipType = zipType
}

type NihaoAnalysis struct{
	Active       string `json:"active"`
	DpvCmra      string `json:"dpv_cmra"`
	DpvFootnotes string `json:"dpv_footnotes"`
	DpvMatchCode string `json:"dpv_match_code"`
	DpvVacant    string `json:"dpv_vacant"`
}
func (this *NihaoAnalysis) GetActive() string {
	if this == nil {
		return ""
	}
	return this.Active
}
func (this *NihaoAnalysis) SetActive(active string) {
	if this == nil {
		return
	}
	this.Active = active
}
func (this *NihaoAnalysis) GetDpvCmra() string {
	if this == nil {
		return ""
	}
	return this.DpvCmra
}
func (this *NihaoAnalysis) SetDpvCmra(dpvCmra string) {
	if this == nil {
		return
	}
	this.DpvCmra = dpvCmra
}
func (this *NihaoAnalysis) GetDpvFootnotes() string {
	if this == nil {
		return ""
	}
	return this.DpvFootnotes
}
func (this *NihaoAnalysis) SetDpvFootnotes(dpvFootnotes string) {
	if this == nil {
		return
	}
	this.DpvFootnotes = dpvFootnotes
}
func (this *NihaoAnalysis) GetDpvMatchCode() string {
	if this == nil {
		return ""
	}
	return this.DpvMatchCode
}
func (this *NihaoAnalysis) SetDpvMatchCode(dpvMatchCode string) {
	if this == nil {
		return
	}
	this.DpvMatchCode = dpvMatchCode
}
func (this *NihaoAnalysis) GetDpvVacant() string {
	if this == nil {
		return ""
	}
	return this.DpvVacant
}
func (this *NihaoAnalysis) SetDpvVacant(dpvVacant string) {
	if this == nil {
		return
	}
	this.DpvVacant = dpvVacant
}

type NihaoComponents struct{
	CityName                string `json:"city_name"`
	DeliveryPoint           string `json:"delivery_point"`
	DeliveryPointCheckDigit string `json:"delivery_point_check_digit"`
	Plus4Code               string `json:"plus4_code"`
	PrimaryNumber           string `json:"primary_number"`
	StateAbbreviation       string `json:"state_abbreviation"`
	StreetName              string `json:"street_name"`
	StreetPredirection      string `json:"street_predirection"`
	StreetSuffix            string `json:"street_suffix"`
	Zipcode                 string `json:"zipcode"`
}
func (this *NihaoComponents) GetCityName() string {
	if this == nil {
		return ""
	}
	return this.CityName
}
func (this *NihaoComponents) SetCityName(cityName string) {
	if this == nil {
		return
	}
	this.CityName = cityName
}
func (this *NihaoComponents) GetDeliveryPoint() string {
	if this == nil {
		return ""
	}
	return this.DeliveryPoint
}
func (this *NihaoComponents) SetDeliveryPoint(deliveryPoint string) {
	if this == nil {
		return
	}
	this.DeliveryPoint = deliveryPoint
}
func (this *NihaoComponents) GetDeliveryPointCheckDigit() string {
	if this == nil {
		return ""
	}
	return this.DeliveryPointCheckDigit
}
func (this *NihaoComponents) SetDeliveryPointCheckDigit(deliveryPointCheckDigit string) {
	if this == nil {
		return
	}
	this.DeliveryPointCheckDigit = deliveryPointCheckDigit
}
func (this *NihaoComponents) GetPlus4Code() string {
	if this == nil {
		return ""
	}
	return this.Plus4Code
}
func (this *NihaoComponents) SetPlus4Code(plus4Code string) {
	if this == nil {
		return
	}
	this.Plus4Code = plus4Code
}
func (this *NihaoComponents) GetPrimaryNumber() string {
	if this == nil {
		return ""
	}
	return this.PrimaryNumber
}
func (this *NihaoComponents) SetPrimaryNumber(primaryNumber string) {
	if this == nil {
		return
	}
	this.PrimaryNumber = primaryNumber
}
func (this *NihaoComponents) GetStateAbbreviation() string {
	if this == nil {
		return ""
	}
	return this.StateAbbreviation
}
func (this *NihaoComponents) SetStateAbbreviation(stateAbbreviation string) {
	if this == nil {
		return
	}
	this.StateAbbreviation = stateAbbreviation
}
func (this *NihaoComponents) GetStreetName() string {
	if this == nil {
		return ""
	}
	return this.StreetName
}
func (this *NihaoComponents) SetStreetName(streetName string) {
	if this == nil {
		return
	}
	this.StreetName = streetName
}
func (this *NihaoComponents) GetStreetPredirection() string {
	if this == nil {
		return ""
	}
	return this.StreetPredirection
}
func (this *NihaoComponents) SetStreetPredirection(streetPredirection string) {
	if this == nil {
		return
	}
	this.StreetPredirection = streetPredirection
}
func (this *NihaoComponents) GetStreetSuffix() string {
	if this == nil {
		return ""
	}
	return this.StreetSuffix
}
func (this *NihaoComponents) SetStreetSuffix(streetSuffix string) {
	if this == nil {
		return
	}
	this.StreetSuffix = streetSuffix
}
func (this *NihaoComponents) GetZipcode() string {
	if this == nil {
		return ""
	}
	return this.Zipcode
}
func (this *NihaoComponents) SetZipcode(zipcode string) {
	if this == nil {
		return
	}
	this.Zipcode = zipcode
}

type Nihao struct{
	Analysis             []*NihaoAnalysis `json:"analysis"`
	CandidateIndex       int32            `json:"candidate_index"`
	Components           *NihaoComponents `json:"components"`
	DeliveryLine_1       string           `json:"delivery_line_1"`
	DeliveryPointBarcode string           `json:"delivery_point_barcode"`
	InputIndex           int32            `json:"input_index"`
	LastLine             string           `json:"last_line"`
	Metadata             *NihaoMetadata   `json:"metadata"`
}
func (this *Nihao) GetAnalysis() []*NihaoAnalysis {
	if this == nil {
		return nil
	}
	return this.Analysis
}
func (this *Nihao) SetAnalysis(analysis []*NihaoAnalysis) {
	if this == nil {
		return
	}
	this.Analysis = analysis
}
func (this *Nihao) GetCandidateIndex() int32 {
	if this == nil {
		return 0
	}
	return this.CandidateIndex
}
func (this *Nihao) SetCandidateIndex(candidateIndex int32) {
	if this == nil {
		return
	}
	this.CandidateIndex = candidateIndex
}
func (this *Nihao) GetComponents() *NihaoComponents {
	if this == nil {
		return nil
	}
	return this.Components
}
func (this *Nihao) SetComponents(components *NihaoComponents) {
	if this == nil {
		return
	}
	this.Components = components
}
func (this *Nihao) GetDeliveryLine_1() string {
	if this == nil {
		return ""
	}
	return this.DeliveryLine_1
}
func (this *Nihao) SetDeliveryLine_1(deliveryLine_1 string) {
	if this == nil {
		return
	}
	this.DeliveryLine_1 = deliveryLine_1
}
func (this *Nihao) GetDeliveryPointBarcode() string {
	if this == nil {
		return ""
	}
	return this.DeliveryPointBarcode
}
func (this *Nihao) SetDeliveryPointBarcode(deliveryPointBarcode string) {
	if this == nil {
		return
	}
	this.DeliveryPointBarcode = deliveryPointBarcode
}
func (this *Nihao) GetInputIndex() int32 {
	if this == nil {
		return 0
	}
	return this.InputIndex
}
func (this *Nihao) SetInputIndex(inputIndex int32) {
	if this == nil {
		return
	}
	this.InputIndex = inputIndex
}
func (this *Nihao) GetLastLine() string {
	if this == nil {
		return ""
	}
	return this.LastLine
}
func (this *Nihao) SetLastLine(lastLine string) {
	if this == nil {
		return
	}
	this.LastLine = lastLine
}
func (this *Nihao) GetMetadata() *NihaoMetadata {
	if this == nil {
		return nil
	}
	return this.Metadata
}
func (this *Nihao) SetMetadata(metadata *NihaoMetadata) {
	if this == nil {
		return
	}
	this.Metadata = metadata
}