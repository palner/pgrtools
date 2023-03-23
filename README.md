# pgrtools

Company Website: <https://www.pgpx.io>

## Overview

pgrtools (aka Palner Group tools) are designed to help with commonly needed functions in go (golang).

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

Man License: GPLv2.

### License / Warranty

The **pgrtools** code is provided in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

pgrtools is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 2 of the License, or (at your option) any later version.

pgrtools is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

## Contact

The Palner Group, Inc.  
Phone: +1 (212) 937-7844  
Matrix: [#help:matrix.lod.com](https://matrix.to/#/#help:matrix.lod.com)
