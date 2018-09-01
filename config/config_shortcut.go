package config

//======== shortcuts for single-value key ==============

// return raw value by key, return error if key not found
// return error if request failed (http driver)
func (c *Config) Get(key string) ([]byte, error) {
	b, err := c.Driver.Get(key)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// return string
func (c *Config) String(key string) string {
	return c.StringOr(key, "")
}

// return string by name, or default value if not found.
func (c *Config) StringOr(key string, dv string) string {
	b, err := c.Get(key)
	if err != nil {
		return dv
	}

	return toString(b, "")
}

func (c *Config) Bool(name string) bool {
	return c.BoolOr(name, false)
}

func (c *Config) BoolOr(name string, dv bool) bool {
	s, err := c.Driver.Get(name)
	if err != nil {
		return dv
	}

	return toBool(s, dv)
}

func (c *Config) Int(name string) int {
	return c.IntOr(name, 0)
}

func (c *Config) IntOr(name string, dv int) int {
	b, err := c.Get(name)
	if err != nil {
		return dv
	}

	return toInt(b, dv)
}

func (c *Config) Int64(name string) int64 {
	return c.Int64Or(name, 0)
}

func (c *Config) Int64Or(name string, dv int64) int64 {
	b, err := c.Get(name)
	if err != nil {
		return dv
	}

	return toInt64(b, dv)
}

func (c *Config) Float64(name string) float64 {
	return c.Float64Or(name, 0)
}

func (c *Config) Float64Or(name string, dv float64) float64 {
	b, err := c.Get(name)
	if err != nil {
		return dv
	}

	return toFloat64(b, dv)
}

//TODO, more types
