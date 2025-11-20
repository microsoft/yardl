# Configurable Array Implementation

By default, Yardl uses the [xtensor](https://xtensor.readthedocs.io/en/latest/) library to
implement multi-dimensional arrays in C++.

Yardl also supports user-defined multi-dimensional array implementations.

If, for example, the `xtensor` types are incompatible with your target software environment, you
can tell Yardl where to find your own implementation at model generation time.

## Defining a Custom Array Implementation

Yardl requires an implementation for each of its three types of multi-dimensional arrays:

1. `FixedNDArray<T, Dims...>`: All dimension sizes are known at compile time
2. `NDArray<T, N>`: Number of dimensions is known at compile time
3. `DynamicNDArray<T>`: Number of dimensions configured at runtime

See [Arrays](language#arrays) for more on the Yardl array types.

These definitions must be accessible from a single "override" header.

### The Override Header

The array override header has three responsibilities:

1. Include your custom array implementation, which may be defined in other C++ headers/source files.
1. Define the Yardl array types.
1. Define the free functions required for compatibility with Yardl.

For example, an override header may look like the following:

```cpp
/** Include multi-dimensional array implementation(s) **/
#include <external-ndarray-implementation>
#include <custom-dynamic-array>
#include <xtensor/containers/xfixed.hpp>


namespace yardl {

/** Define the three array types **/

// Alias xtensor's fixed array type
template <typename T, size_t... Dims>
using FixedNDArray = xt::xtensor_fixed<T, xt::xshape<Dims...>, xt::layout_type::row_major, false>;

// Alias my custom dynamic array implementation
template <typename T>
using DynamicNDArray = custom::dynamic_array<T>;

// Extend my external Array implementation to implement NDArray
template <typename T, size_t N>
class NDArray : public external::Array<T, N>
{
    // ... Wrapper implementation
}

/** API functions required for Yardl compatibility **/

template <typename T, size_t... Dims>
size_t size(FixedNDArray<T, Dims...> const& arr) {
    return arr.size();
}

template <typename T>
size_t size(DynamicNDArray<T> const& arr) {
    return arr.get_size();
}

template <typename T, size_t N>
size_t size(NDArray<T, N> const& arr) {
    return arr.get_number_of_elements();
}

/** More API functions continued... **/

}
```

### General Requirements

1. Your header must define the three array types, which have the following type signatures
    1. `template <typename T, size_t... Dims> class FixedNDArray`
    2. `template <typename T> class DynamicNDArray`
    3. `template <typename T, size_t N> class NDArray`
1. Each type must have:
    1. Default constructor (can be constructed with no arguments)
    1. Copy/Move constructors
    1. Copy/Move assignment operators
    1. Equality/Inequality operators
    1. `begin()` and `end()` member functions to support iteration (`const` and non-`const`)
1. The `FixedNDArray` type must be [Trivally Copyable](https://en.cppreference.com/w/cpp/named_req/TriviallyCopyable). For this reason, using `xtensor` for fixed arrays is recommended.
1. The array types must be defined within the `yardl` namespace
1. Your header must also define a collection of API functions, described below, also within the `yardl` namespace


### Array API Requirements

The following functions must be defined for *each* array type `A<T>`:

1. `size_t size(A const& a)`

    Returns the total number of elements.

1. `size_t dimension(A const& a)`

    Returns the number of dimensions.

1. `size_t shape(A const& a, size_t dimension)`

    Returns the length of the given dimension.

1. `T* dataptr(A<T>& a)`

    Returns a pointer to the first element.

1. `T const* dataptr(A<T> const& a)`

    Returns a `const` pointer to the first element.

1. `template <typename T, class... Args> T const& at(A<T> const& a, Args... indices)`

    Returns a `const` reference to the element at the given indices.


The following functions must be defined in addition to those above:

1. `std::array<size_, sizeof...(Dims)> shape(FixedNDArray<T, Dims...> const& a)`

    Returns the shape (dimensions) of the fixed array.

1. `std::vector<size_t> shape(DynamicNDArray<T> const& a)`

    Returns the shape (dimensions) of the dynamic array.

1. `std::array<size_t> shape(NDArray<T, N> const& a)`

    Returns the shape (dimensions) of the ndarray.

1. `void resize(DynamicNDArray<T>& a, std::vector<size_t> const& shape)`

    Changes the array's shape (dimensions) without preserving data.

1. `void resize(NDArray<T, N>& a, std::array<size_t> const& shape)`

    Changes the array's shape (dimensions) without preserving data.


## Using a Custom Array Implementation

To configure Yardl to use your custom array implementation in your model, add the `overrideArrayHeader`
option to the `cpp` section of your model's `_package.yml`.

For example:

```yaml
namespace: MyNamespace

cpp:
  sourcesOutputDir: ../cpp/generated
  generateCMakeLists: true
  overrideArrayHeader: external/my-array-impl.h     // [!code ++]
```

The generated C++ code will then `#include "external/my-array-impl.h"` instead of the default Yardl
multi-dimensional array implementation.

Ensure that your header is on the include path before compiling your generated code.


## Examples

Yardl provides two examples for configuring multi-dimensional array implementations.

### Yardl's Default Implementation

Yardl's default `xtensor` implementation can be found in `yardl/detail/ndarray/impl.h` after
generating C++ code for your model.

This header:
1. Includes the `xtensor` headers
2. Aliases the `xtensor` types to define each of the three Yardl array types
3. Implements *all* of the API functions described above

### Example Custom Implementation

This example can be found in the Yardl code in the `cpp/test/external` directory.
This directory contains two files:

1. `hoNDArray.h`

    This file defines a base `NDArray<T>` class and derived `hoNDArray<T>` class, both in the
    `external` namespace. These are loosely based on the multi-dimensional array implementation
    found in the Gadgetron open source image reconstruction framework.

2. `ndarray_impl.h`

    This file:
    1. Includes `hoNDArray.h`
    1. Defines a collection of helpers in the `yardl::detail` namespace
        - These are used to implement syntactic convenience and may not be needed in your code
    1. Uses `xtensor` to implement `yardl::FixedNDArray` (identical to Yardl's default implementation)
    1. Defines two *new* classes `yardl::NDArray` and `yardl::DynamicNDArray`, both of which derive
        from the `external::hoNDArray` implementation.
    1. Defines the array API free functions
        -  Since both `NDArray` and `DynamicNDArray` extend the same `hoNDArray` class, many of the
            free functions are only implemented once, e.g. the following function implements the
            `yardl::size` operator for both those types:
            ```cpp
            template <typename T>
            size_t size(external::hoNDArray<T> const& arr);
            ```
