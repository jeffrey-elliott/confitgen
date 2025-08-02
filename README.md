# confitgen

confitgen is a lightweight Go code generator for schema-driven config types.  
Given a JSON schema, it emits typed structs and a tiny bit of boilerplate for loading values.

## Usage

```
go run . -schema <package>
```

This expects a schema file like:

```
<package>.confit.schema.json
```

For example:

```
go run . -schema starship
```

Looks for:

```
starship.confit.schema.json
```

Schema input (`starship.confit.schema.json`):

```json
{
  "StarshipConfitSchema": {
    "Name": "string",
    "Class": "string",
    "Drive": "string",
    "InfiniteImprobabilityEnabled": "bool",
    "Destination": "string",
    "Crew": "[]string"
  }
}
```

It generates:

```
starship.go
```

Generated output:

```go
type Starship struct {
	Destination                  string   `json:"Destination"`
	Crew                         []string `json:"Crew"`
	Name                         string   `json:"Name"`
	Class                        string   `json:"Class"`
	Drive                        string   `json:"Drive"`
	InfiniteImprobabilityEnabled bool     `json:"InfiniteImprobabilityEnabled"`
}
```

Place that generated file in your own app:

```
internal/confit/starship.go
```

## Why did you write this? 

Config turns up everywhere as a problem with lots of different solutions. Most of the time you either buy into someone else's solution (far too much, and another dependency) or roll your own (probably too little, and lots of repeated effort). confitgen is my way to have a declarative configuration I can drop into another project without adding yet another dependency.

## Why is this named confitgen?

Confit preserves your config in the same way duck confit preserves the duck, leaving you with delicious, self-contained config morsels ready to drop into any dish... er, project. Slow-cooked type safety, no grease. No extra dependencies.
