package record

type ItemRecord struct {
	Category    string
	Item        string
	Description string
}

type ProductRecord struct {
	Category              string
	Item                  string
	Product               string
	PartNumber            string
	AcceptBeforeSixMonths bool
	AcceptAfterSixMonths  bool
	Price                 int
}
