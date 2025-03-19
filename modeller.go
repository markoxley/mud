package dtorm

import (
	"reflect"
	"strconv"
	"strings"
)

type Modeller interface {
	// StandingData returns the standing data for the model
	StandingData() []Modeller

	// GetID returns the ID of the model
	GetID() *string

	// IsNew returns true if the model has yet to be saved
	IsNew() bool

	// IsDeleted returns true if the model has been marked as deleted
	IsDeleted() bool

	Disable()
}

func getDefs(t interface{}, first bool) []field {

	res := make([]field, 0, 10)
	if first {
		res = append(res, field{
			name:     "ID",
			fType:    tUUID,
			identity: true,
			key:      true,
		})
		res = append(res, field{
			name:  "CreateDate",
			fType: tDateTime,
			key:   true,
		})
		res = append(res, field{
			name:  "LastUpdate",
			fType: tDateTime,
			key:   true,
		})
		res = append(res, field{
			name:      "DeleteDate",
			fType:     tDateTime,
			allowNull: true,
		})
	}
	v := reflect.ValueOf(t)
	ft := reflect.TypeOf(t)
	nf := v.NumField()
	for i := 0; i < nf; i++ {

		st := ft.Field(i)
		sv := v.Field(i)
		null := false
		if sv.Kind() == reflect.Ptr {
			null = true
			sv = sv.Elem()
		}

		if sv.Kind() == reflect.Struct && sv.Type().Name() != "Time" {
			if subf := getDefs(sv.Interface(), false); len(subf) > 0 {
				res = append(res, subf...)
			}
		} else {
			if tg, ok := st.Tag.Lookup("mxorm"); ok {
				nm := st.Name
				szMj := 0
				szMn := 0
				id := false
				key := false
				uns := false
				fld := tString

			FieldSearchLoop:
				for k, v := range fieldTrans {
					for _, v2 := range v {
						if v2 == sv.Kind() {
							fld = k
							for _, sn := range fieldUnsigned {
								if sn == sv.Kind() {
									uns = true
								}
							}
							break FieldSearchLoop
						}
					}
				}
				if tg != "" {
					tgs := strings.Split(tg, ",")
					for _, t := range tgs {
						pts := strings.Split(t, ":")
						if len(pts) == 2 {
							switch pts[0] {
							case "type":
								typeKey := pts[1]
								if typeKey == "time" {
									typeKey = "struct"
								}
								if v, ok := fieldNames[typeKey]; ok {

									fld = v
								}
							case "size":
								szPt := strings.Split(pts[1], ",")
								if v, err := strconv.ParseInt(szPt[0], 10, 64); err == nil {
									szMj = int(v)
									if len(szPt) > 1 {
										if v, err = strconv.ParseInt(szPt[1], 10, 64); err == nil {
											szMn = int(v)
										}
									}
								}
							case "identity":
								id = pts[1] == "true"
							case "key":
								key = pts[1] == "true"
							case "unsigned":
								uns = pts[1] == "true"
							}
						}
					}
				}
				res = append(res, newField(nm, fld, szMj, szMn, id, key, uns, null))
			}
		}
	}
	return res
}
