# pgrtools

Company Website: [The Palner Group, Inc.](https://www.palner.com)

## Overview

pgrtools (**P**alner **Gr**oup **Tools**) are designed to help with commonly needed functions in go (golang).

## Usage Example

```go
...
	import (
		...
		"github.com/palner/pgrtools/pgparse"
		"github.com/palner/pgrtools/pgsqljson"
		...
	)
...
```

```go
...
	requiredfields := []string{"ipaddress", "shortnote"}
	keyVal, err := pgparse.ParseBodyFields(r, requiredfields)
...
```

## License 

License: MIT

### License / Warranty

**pgrtools** code is provided in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

pgrtools is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the MIT for more details.

## Contact

The Palner Group, Inc.  
Phone: +1 (212) 937-7844  
Matrix: [#help:matrix.lod.com](https://matrix.to/#/#help:matrix.lod.com)
