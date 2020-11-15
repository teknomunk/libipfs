
# Library Methodology

The interface was built to follow the following rules:

* Handles for all types that are not
  * integer or floating point numbers
  * strings
  * booleans
  * enumeration
  * callback of the type  ``int(*)( void* /* callback_data */, /* extra arguments */ )``
* Separate function calls for setting/getting fields in complex types
* standard prefix for function calls ipfs_

