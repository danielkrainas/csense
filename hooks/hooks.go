package hooks

import (
	"regexp"
	"time"

	"github.com/danielkrainas/gobag/util/uuid"

	"github.com/danielkrainas/csense/api/v1"
)

type Filter interface {
	Match(hook *v1.Hook, c *v1.ContainerInfo) bool
}

type CriteriaFilter struct{}

func (f *CriteriaFilter) Match(hook *v1.Hook, c *v1.ContainerInfo) bool {
	crit := hook.Criteria

	for fieldName, condition := range crit.Fields {
		valid := false
		switch fieldName {
		case v1.FieldName:
			valid = IsValid(condition, c.Name)
		case v1.FieldImageName:
			valid = IsValid(condition, c.ImageName)
		}

		if valid {
			return valid
		}
	}

	for k, v := range c.Labels {
		if x, ok := c.Labels[k]; ok && x == v {
			return true
		}
	}

	return false
}

func IsValid(c *v1.Condition, v string) bool {
	if c == nil {
		return false
	}

	switch c.Op {
	case v1.OperandEqualShort:
		fallthrough
	case v1.OperandEqual:
		return c.Value == v

	case v1.OperandNotEqualShort:
		fallthrough
	case v1.OperandNotEqual:
		return c.Value != v

	case v1.OperandMatch:
		ok, err := regexp.MatchString(c.Value, v)
		return err == nil && ok
	}

	return false
}

func DefaultHook() *v1.Hook {
	return &v1.Hook{
		ID:      uuid.Generate(),
		Events:  make([]v1.EventType, 0),
		TTL:     -1,
		Created: time.Now().Unix(),
		Format:  v1.FormatJSON,
	}
}

func FilterAll(hooks []*v1.Hook, c *v1.ContainerInfo, f Filter) []*v1.Hook {
	results := make([]*v1.Hook, 0)
	for _, hook := range hooks {
		if f.Match(hook, c) {
			results = append(results, hook)
		}
	}

	return results
}
