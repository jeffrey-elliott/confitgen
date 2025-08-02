func New() TypeName { return TypeName{} }
func Load(path string) (TypeName, error) {
cf := New()
data, err := os.ReadFile(path)
if err != nil { return cf, err }
if err := json.Unmarshal(data, &cf); err != nil { return cf, err }
return cf, nil
}
