package value_render

import (
	"fmt"
	"reflect"
	"regexp"
	"time"
)

func dateFormat(t interface{}, format string) (string, error) {
	if reflect.TypeOf(t).String() == "time.Time" {
		t1 := t.(time.Time)
		return t1.Format(format), nil
	}
	if reflect.TypeOf(t).Kind() == reflect.Int {
		t1 := int64(t.(int))
		return time.Unix(t1/1000, t1%1000*1000000).Format(format), nil
	}
	if reflect.TypeOf(t).Kind() == reflect.Int64 {
		t1 := t.(int64)
		return time.Unix(t1/1000, t1%1000*1000000).Format(format), nil
	}
	if reflect.TypeOf(t).Kind() == reflect.String {
		t1, e := time.Parse(time.RFC3339, t.(string))
		if e != nil {
			return format, e
		}
		//fmt.Println(t1)
		return t1.Format(format), nil
	}
	return format, fmt.Errorf("could not tell the type timestamp field belongs to")
}

type IndexRender struct {
	dateFormat  string
	valueFormat string
}

func NewIndexRender(t string) *IndexRender {
	r, _ := regexp.Compile(`%{\+.*?}`)
	loc := r.FindStringIndex(t)
	return &IndexRender{
		dateFormat:  t[loc[0]+2 : loc[1]-1],
		valueFormat: t[:loc[0]] + "%s" + t[loc[1]+1:],
	}
}

func (r *IndexRender) Render(event map[string]interface{}) interface{} {
	var s string
	if t, ok := event["@timestamp"]; ok {
		s, _ = dateFormat(t, r.dateFormat)
	} else {
		s, _ = dateFormat(time.Now(), r.dateFormat)
	}
	return fmt.Sprintf(r.valueFormat, s)
}
