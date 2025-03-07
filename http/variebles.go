package http

type Varibles map[string]string

func (v *Varibles) Add(key string, value string) {
	(*v)[key] = value
}
func (v *Varibles) Get(key string) string {
	return (*v)[key]
}
func (v *Varibles) Set(key string, value string) {
	(*v)[key] = value
}
func (v *Varibles) Del(key string) {
	delete((*v), key)
}
func (v *Varibles) Clear() {
	(*v) = make(map[string]string)
}
func (v *Varibles) Count() int {
	return len((*v))
}
func (v *Varibles) Exist(key string) bool {
	_, ok := (*v)[key]
	return ok
}
func (v *Varibles) Keys() []string {
	keys := make([]string, 0, len((*v)))
	for k := range *v {
		keys = append(keys, k)
	}
	return keys
}
func (v *Varibles) Values() []string {
	values := make([]string, 0, len((*v)))
	for _, v := range *v {
		values = append(values, v)
	}
	return values
}
func (v *Varibles) ToMap() map[string]string {
	return (*v)
}
func NewVaribles() Varibles {
	return make(Varibles)
}
