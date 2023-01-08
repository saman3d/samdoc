package docx

import (
	"os"
	"testing"
)

type populate_data struct {
	fullname      string
	gender        string
	contract_code string
	contract_date string
	seller_name   string
	father_name   string
	national_no   string
	identity_no   string
	issuing_city  string
	bith_date     string
	birth_city    string
	r             struct {
		plaque      string
		alley       string
		street      string
		region      string
		city        string
		province    string
		phone_no    string
		postal_code string
	}
	mobile_no               string
	email                   string
	car_brand               string
	car_model               string
	car_series              string
	total_amount_in_words   string
	total_amount_in_digit   string
	prepayment              string
	final_debit_balance     string
	installment_count       string
	installment_amount      string
	delivery_time           string
	cheque_no               string
	bank_neme               string
	cheque_amount           string
	cheque_date             string
	payable_account         string
	guarantor_fulname       string
	guarantor_gender        string
	guarantor_national_code string
}

var data = populate_data{
	fullname:      "سامان کوشکی",
	gender:        "آقای",
	contract_code: "2323233",
	contract_date: "1401/10/02",
	seller_name:   "علی اصغر",
	father_name:   "حسن",
	national_no:   "2323",
	identity_no:   "233",
	issuing_city:  "سبزوار",
	bith_date:     "1379/01/02",
	birth_city:    "سبزوار",
	r: struct {
		plaque      string
		alley       string
		street      string
		region      string
		city        string
		province    string
		phone_no    string
		postal_code string
	}{
		plaque:      "تهران",
		alley:       "عربعلی",
		street:      "عربعلی",
		region:      "2",
		city:        "تهران",
		province:    "تهران",
		phone_no:    "+989331736324",
		postal_code: "34234234",
	},
	mobile_no:               "+989331736324",
	email:                   "saman.koushki79@gmail.com",
	car_brand:               "toyota",
	car_model:               "x3",
	car_series:              "ss",
	total_amount_in_words:   "سی۲د ملیون و هفتصد",
	total_amount_in_digit:   "123123",
	prepayment:              "20000",
	final_debit_balance:     "4000000",
	installment_count:       "24",
	installment_amount:      "50000",
	delivery_time:           "۴۰ روز",
	cheque_no:               "1000",
	bank_neme:               "سامان",
	cheque_amount:           "2323",
	cheque_date:             "1401/02/03",
	payable_account:         "34234",
	guarantor_fulname:       "اصغر فرهادی",
	guarantor_gender:        "آقا",
	guarantor_national_code: "3423234234",
}

func TestTemplate(t *testing.T) {
	finp, err := os.Open("test.docx")
	if err != nil {
		panic(err)
	}

	tmp, err := NewTemplate(finp)
	if err != nil {
		panic(err)
	}

	fout, err := os.Create("test_out.docx")
	if err != nil {
		panic(err)
	}

	err = tmp.ExecuteToWriter(&data, fout)
	if err != nil {
		panic(err)
	}

}

func TestTemplatePDF(t *testing.T) {
	finp, err := os.Open("test.docx")
	if err != nil {
		panic(err)
	}

	tmp, err := NewTemplate(finp)
	if err != nil {
		panic(err)
	}

	fout, err := os.Create("test_out.pdf")
	if err != nil {
		panic(err)
	}

	pdf, err := tmp.ExecuteToPDF(&data)
	if err != nil {
		panic(err)
	}

	fout.Write(pdf)
}
