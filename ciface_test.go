package ciface_test

import (
	. "github.com/npolar/ciface"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ciface", func() {
	var (
		simpleCsv    CsvInterface
		manualHeader CsvInterface
		hiragana     *CsvInterface
	)

	BeforeEach(func() {
		simpleCsv = CsvInterface{
			Data: []byte("latitude,longitude,temp,dtime\n69.1833,18.5632,-5.3,2015-01-02T10:15:12Z"),
		}

		manualHeader = CsvInterface{
			Data:   []byte("83.67,17.12,Lance,true,278.3"),
			Header: []string{"lat", "lng", "name", "online", "course"},
		}

	})

	Describe("Standard operation", func() {
		Context("Initialization", func() {
			It("Should load []byte data", func() {
				hiragana = NewParser([]byte("#First set of hiragana characters\nあ;い;う;え;お"))
				Expect(hiragana.Data).To(Equal([]byte("#First set of hiragana characters\nあ;い;う;え;お")))
			})

			It("Set custom header", func() {
				hiragana.Header = []string{"a", "i", "u", "e", "o"}
				Expect(hiragana.Header).To(Equal([]string{"a", "i", "u", "e", "o"}))
			})

			It("Set custom delimiter", func() {
				hiragana.Delimiter = ';'
				Expect(hiragana.Delimiter).To(Equal(';'))
			})

			It("Set comment rune", func() {
				hiragana.Comment = '#'
				Expect(hiragana.Comment).To(Equal('#'))
			})
		})

		Context("Parsing", func() {
			It("should return あ for a key", func() {
				hiragana = NewParser([]byte("a,i,u,e,o\nあ,い,う,え,お"))
				a := hiragana.Parse()
				Expect(a[0].(map[string]interface{})["a"]).To(Equal("あ"))
			})

			It("should ignore the comment", func() {
				hiragana = NewParser([]byte("#my random comment\na,i,u,e,o\nあ,い,う,え,お"))
				hiragana.Comment = '#'
				a := hiragana.Parse()
				Expect(a[0].(map[string]interface{})["a"]).To(Equal("あ"))
			})

			It("should split on ; when defined as delimiter", func() {
				hiragana = NewParser([]byte("a;i;u;e;o\nあ;い;う;え;お"))
				hiragana.Delimiter = ';'
				a := hiragana.Parse()
				Expect(a[0].(map[string]interface{})["o"]).To(Equal("お"))
			})

			It("should handle winblows new lines", func() {
				hiragana = NewParser([]byte("a;i;u;e;o\r\nあ;い;う;え;お"))
				hiragana.Delimiter = ';'
				a := hiragana.Parse()
				Expect(a[0].(map[string]interface{})["e"]).To(Equal("え"))
			})

			It("should cast true strings to a boolean", func() {
				kana := NewParser([]byte("hiragana,katakana,romaji,symbol\ntrue,false,na,な"))
				a := kana.Parse()
				Expect(a[0].(map[string]interface{})["hiragana"]).To(Equal(true))
			})

			It("should cast false strings to a boolean", func() {
				kana := NewParser([]byte("hiragana,katakana,romaji,symbol\ntrue,FALSE,na,な"))
				a := kana.Parse()
				Expect(a[0].(map[string]interface{})["katakana"]).To(Equal(false))
			})

			It("should parse numeric strings as numbers", func() {
				数字 := NewParser([]byte("symbol,number\n一,1"))
				a := 数字.Parse()
				Expect(a[0].(map[string]interface{})["number"]).To(Equal(float64(1)))
			})

			It("should round floats to the proper precision", func() {
				number := NewParser([]byte("latitude,longitude\n69.21342344,18.234512234"))
				number.Precision = 5
				a := number.Parse()
				Expect(a[0].(map[string]interface{})["latitude"]).To(Equal(69.21342))
			})
		})
	})

	Describe("Helper Functions", func() {

		Context("Boolean string detection", func() {
			It("should be true when getting TRuE string", func() {
				Expect(BooleanString("true")).To(Equal(true))
			})

			It("should be true when getting fAlse string", func() {
				Expect(BooleanString("false")).To(Equal(true))
			})

			It("should be false when getting random string", func() {
				Expect(BooleanString("何")).To(Equal(false))
			})

			It("should be false when getting 1 string", func() {
				Expect(BooleanString("1")).To(Equal(false))
			})
		})

		Context("Rounding correction", func() {
			It("return 4.444 with precision 3 and input 4.4443323", func() {
				Expect(Round(4.4443323, 3)).To(Equal(4.444))
			})

			It("return 4.45 with precision 2 and input 4.4474324", func() {
				Expect(Round(4.4474324, 2)).To(Equal(4.45))
			})

			It("return 5.0 with precision 1 and input 4.99999", func() {
				Expect(Round(4.99999, 1)).To(Equal(5.0))
			})
		})
	})
})
