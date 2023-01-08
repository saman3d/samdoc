package samdoc

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldQuery(t *testing.T) {
	Convey("Test Field Query: Pop", t, func() {
		fq := FieldQuery([]string{"saman", "ali", "asghar", "mojtaba"})

		fq = fq.Pop(0)
		So(fq, ShouldResemble, FieldQuery([]string{"ali", "asghar", "mojtaba"}))

		fq = fq.Pop(1)
		So(fq, ShouldResemble, FieldQuery([]string{"ali", "mojtaba"}))

		fq = fq.Pop(0)
		So(fq, ShouldResemble, FieldQuery([]string{"mojtaba"}))

		fq = fq.Pop(0)
		So(fq, ShouldResemble, FieldQuery([]string{}))
	})
}

type NIN struct {
	f string
	s string
}

func (n NIN) String() string {
	return n.f + n.s
}

type Cert struct {
	Born time.Time
	NIN   NIN
	Image string
}

type Address struct {
	Country   string
	State     string
	City      string
	Street    string
	PlaqueNum int
}

func (a Address) String() string {
	return fmt.Sprintf("%s, %s, %s, %s, plaque %d", a.Country, a.State, a.City, a.Street, a.PlaqueNum)
}

type Person struct {
	Name     string
	LastName string
	Cert     *Cert
	Address  Address
	Weight int
}

func TestStructure(t *testing.T) {
	Convey("Test Structure: Get Field", t, func() {

		var saman = Person{
			Name:     "saman",
			LastName: "koushki",
			Cert: &Cert{
				NIN: NIN{
					f: "1",
					s: "2",
				},
				Image: "0280374",
			},
		}

		m, err := NewStructure(saman)

		So(err, ShouldBeNil)
		nin, err := m.Get(FieldQuery{
			"Cert",
			"NIN",
		})
		So(err, ShouldBeNil)
		So(nin, ShouldEqual, saman.Cert.NIN.String())

	})
}
