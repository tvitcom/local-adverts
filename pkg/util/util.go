package util

import (
	"strconv"
)

type (
	// UserInfo struct {
	// 	User_id  int64
	// 	Name     string
	// 	Picture  string
	// 	Roles    []string
	// 	Phone    string
	// 	Email    string
	// 	Birthday string
	// 	Data     interface{}
	// }
	// PageData struct {
	// 	Name             string
	// 	Lang             string
	// 	BaseUrl          string
	// 	Seo              *Seo
	// 	Biz              *BizInfo
	// 	User             *UserInfo
	// 	LeftMenuSelected int
	// 	Data             interface{}
	// }
	// BizInfo struct {
	// 	Name, ShortName, Email, Phone, Phone2 string
	// }
	// Seo struct {
	// 	Jsonld      string
	// 	Og          *OG
	// 	Description string
	// 	Keywords    string
	// }
	// OG struct { // Open Graph basic struct for SEO
	// 	Title, Description, Type, Url, Image string
	// }
)

// func GetBizInfo() *BizInfo {
// 	biz := new(BizInfo)
// 	biz.ShortName = conf.BIZ_SHORTNAME
// 	biz.Name = conf.BIZ_NAME
// 	biz.Email = conf.BIZ_EMAIL
// 	biz.Phone = conf.BIZ_PHONE
// 	biz.Phone2 = conf.BIZ_PHONE2
// 	return biz
// }

//TODO: make normal jsonld builder
// func GetSeoDefault(c *gin.Context, title, descr, typ, url, logo string) *Seo {
// 	if url == "" {
// 		url = c.FullPath()
// 	}
// 	if typ == "" {
// 		typ = "website"
// 	}
// 	if descr == "" {
// 		descr = "simple useful project application for people"
// 	}
// 	seo := new(Seo)
// 	seo.Jsonld = ""
// 	seo.Og = new(OG)
// 	seo.Og.Title = title
// 	seo.Og.Description = descr
// 	seo.Og.Type = typ
// 	seo.Og.Url = url
// 	seo.Og.Image = logo
// 	seo.Keywords = "simple useful project application for peaple"
// 	seo.Description = descr
// 	return seo
// }

func GetUserpicURL(filepath string) string {
	if filepath == "" {
		filepath = "0_60x60.jpeg"
	}
	return "/userpic/" + filepath
}

func Pkeyer(raw interface{}) int64 {
	var i int64
    switch raw.(type) {
	    case string:
		    str := raw.(string)
		    tp, err := strconv.ParseInt(str, 10, 64)
		    i = tp
			if err != nil {
				panic(err)
			}
	    case int64:
		    i = raw.(int64)
        default:
            i = 0
    }
    return i
}

func ExtractDigitsInt(str string) int {
	if s, err := strconv.ParseInt(ExtractDigitsString(str), 10, 64); err == nil {
	    return int(s)
	}
	return 0
}
