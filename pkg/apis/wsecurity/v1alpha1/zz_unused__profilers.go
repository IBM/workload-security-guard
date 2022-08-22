package v1alpha1

/*
type Set struct {
	List []string
	m    map[string]bool
}

func (p *Set) Add(v string) {
	if !p.m[v] {
		p.m[v] = true
		p.List = append(p.List, v)

	}
}

func (p *Set) Clear() {
	p.m = make(map[string]bool)
	p.List = make([]string, 0, 1)
}

func (p *Set) Merge(a *Set) {
	if p.m == nil {
		p.m = make(map[string]bool, len(a.List))
	}
	if p.List == nil {
		p.List = make([]string, 0, len(a.List))
	}

	for _, v := range a.List {
		if !p.m[v] {
			p.m[v] = true
			p.List = append(p.List, v)
		}
	}
}

func (p *Set) Learn(a *Set) {
	if p.m == nil {
		p.m = make(map[string]bool, len(a.List))
	}
	if p.List == nil {
		p.List = make([]string, 0, len(a.List))
	}

	for _, v := range a.List {
		if !p.m[v] {
			p.m[v] = true
			p.List = append(p.List, v)
		}
	}
}

func (p *Set) Fuse(a *Set) {
	if p.m == nil {
		p.m = make(map[string]bool, len(a.List))
	}
	if p.List == nil {
		p.List = make([]string, 0, len(a.List))
	}

	for _, v := range a.List {
		if !p.m[v] {
			p.m[v] = true
			p.List = append(p.List, v)
		}
	}
}

func (p Set) Decide(s string) string {
	_, exists := p.m[s]
	if !exists {
		return fmt.Sprintf("Unexpected key %s in Set", s)
	}
	return ""
}

func AddToSetFromList(list []string, set *Set) {
	if set.m == nil {
		set.m = make(map[string]bool, len(list))
	}
	if set.List == nil {
		set.List = make([]string, 0, len(list))
	}
	for _, v := range list {
		if v != "" && !set.m[v] {
			set.m[v] = true
			set.List = append(set.List, v)
		}
	}
}

func AddToListFromSet(set *Set, list *[]string) {
	*list = append(*list, set.List...)
}


type Uint32Slice []uint32

func (base Uint32Slice) Add(val Uint32Slice) Uint32Slice {
	if missing := len(val) - len(base); missing > 0 {
		// Dynamically allocate as many blockElements as needed for this config
		base = append(base, make([]uint32, missing)...)
	}
	for i, v := range val {
		base[i] = base[i] | v
	}
	return base
}

func (base Uint32Slice) Decide(val Uint32Slice) string {
	for i, v := range val {
		if v == 0 {
			continue
		}
		if i < len(base) && (v & ^base[i]) == 0 {
			continue
		}
		return fmt.Sprintf("Unexpected Unicode Flag %x on Element %d", v, i)
	}
	return ""
}
func (uint64Slice Uint32Slice) String(depth int) string {
	if uint64Slice == nil {
		return "null"
	}
	if len(uint64Slice) == 0 {
		return "[]"
	}
	var description bytes.Buffer
	description.WriteString(fmt.Sprintf("[%x", uint64Slice[0]))
	for i := 1; i < len(uint64Slice); i++ {
		description.WriteString(fmt.Sprintf(",%x", uint64Slice[i]))
	}
	description.WriteString("]")
	return description.String()
}

*/
/*
func (uint64Slice Uint32Slice) Describe() string {
	if uint64Slice != nil {
		var description bytes.Buffer
		description.WriteString("Unicode slice: ")
		description.WriteString(fmt.Sprintf("%x", uint64Slice[0]))
		for i := 1; i < len(uint64Slice); i++ {
			description.WriteString(fmt.Sprintf(", %x", uint64Slice[i]))
		}
		return description.String()
	}
	return ""
}
*/

/*
func process(o interface{}) {
	if reflect.ValueOf(o).Kind() == reflect.Map {
		fmt.Printf("%v is a map!\n", o)
		for k, v := range o.(map[string]interface{}) {

			fmt.Printf("%s: ", k)
			process(v)
		}
		return
	}
	if reflect.ValueOf(o).Kind() == reflect.Slice {
		fmt.Printf("%v is a slice!\n", o)
		for i, v := range o.([]interface{}) {
			fmt.Printf("%d: ", i)
			process(v)
		}
		return
	}
	fmt.Printf("%v is just a value\n", o)
}
*/

// Profile a PostForm
/* hopefully can be handled with Profile() - It is possible that this is not ok due to the array support limitations
func (profile *StructuredProfile) profilePostForm(m map[string][]string) {
	profile.Kind = PROFILE_OBJECT
	//fmt.Printf("JsonProfile ProfilePostForm\n")
	profile.Kv = make(map[string]*StructuredProfile, len(m))
	for k, sv := range m {
		profile.Kv[k] = new(StructuredProfile)
		profile.Kv[k].Kind = PROFILE_ARRAY
		profile.Kv[k].Vals = make([]SimpleValProfile, len(sv))
		for i, v := range sv {
			profile.Kv[k].Vals[i].Profile(v)
		}
	}
}
*/
