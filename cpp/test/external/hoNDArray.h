#pragma once

#include <vector>

namespace external {

/** Base NDArray class
 *
 *  This class is loosely adapted from the Gadgetron NDArray abstract type.
 *
 *  This implementation expects and stores dimensions in reverse order with
 *  respect to yardl/numpy/xtensor, but the underlying array data is stored
 *  in row-major order.
 */
template <typename T>
class NDArray {
 public:
  NDArray() : data_(0) {};

  void create(std::vector<size_t> const& dimensions) {
    for (size_t dim : dimensions) {
      if (dim == 0) {
        throw std::runtime_error("NDArray::create, dimension size is 0");
      }
    }
    this->dimensions_ = dimensions;
    allocate_memory();
  }

  void clear() {
    this->deallocate_memory();

    this->data_ = nullptr;
    this->dimensions_.clear();
  }

  size_t get_size() const {
    if (this->dimensions_.empty()) {
      return 0;
    }

    size_t s = 1;
    for (size_t dim : this->dimensions_) {
      s *= dim;
    }
    return s;
  }

  inline size_t get_size(size_t dimension) const {
    if (dimension >= this->dimensions_.size()) {
      return 1;
    } else {
      return this->dimensions_[dimension];
    }
  }

  std::vector<size_t> const& get_dimensions() const {
    return this->dimensions_;
  }

  size_t get_number_of_dimensions() const {
    return this->dimensions_.size();
  }

  bool dimensions_equal(NDArray<T> const& a) const {
    return this->dimensions_ == a.dimensions_;
  }

  T* get_data() {
    return this->data_;
  }
  T const* get_data() const {
    return this->data_;
  }

 protected:
  virtual void allocate_memory() = 0;
  virtual void deallocate_memory() = 0;

 protected:
  T* data_;
  std::vector<size_t> dimensions_;
};

/** CPU-based NDArray (Host) */
/** Host (CPU) NDArray class
 *
 *  This class is loosely adapted from the Gadgetron hoNDArray type.
 *
 *  This implementation expects and stores dimensions in reverse order with
 *  respect to yardl/numpy/xtensor, but the underlying array data is stored
 *  in row-major order.
 */
template <typename T>
class hoNDArray : public NDArray<T> {
  using base_type = NDArray<T>;

 public:
  hoNDArray() : base_type() {};

  ~hoNDArray() {
    deallocate_memory();
  }

  // Copy constructor
  hoNDArray(hoNDArray<T> const& a) {
    this->create(a.get_dimensions());
    std::copy(a.get_data(), a.get_data() + a.get_size(), this->get_data());
  }

  // Move constructor
  hoNDArray(hoNDArray<T>&& a) noexcept {
    this->data_ = a.data_;
    this->dimensions_ = a.dimensions_;
    a.data_ = nullptr;
    a.dimensions_.clear();
  }

  // Move assignment operator
  hoNDArray& operator=(hoNDArray&& rhs) noexcept {
    if (&rhs == this)
      return *this;

    this->clear();

    this->data_ = rhs.data_;
    rhs.data_ = nullptr;

    this->dimensions_ = rhs.dimensions_;
    rhs.dimensions_.clear();

    return *this;
  }

  // Copy assignment operator
  hoNDArray& operator=(hoNDArray const& rhs) {
    if (&rhs == this)
      return *this;

    if (rhs.get_size() == 0) {
      this->clear();
      return *this;
    }

    if (!this->dimensions_equal(rhs)) {
      deallocate_memory();
      this->data_ = 0;
      this->dimensions_ = rhs.dimensions_;
      allocate_memory();
    }

    std::copy(rhs.begin(), rhs.end(), this->begin());
    return *this;
  }

  /* Get element at index */
  T const& operator()(std::vector<size_t> const& ind) const {
    // Calculate offset factor for each dimension
    size_t noffsets = this->dimensions_.size();
    std::vector<size_t> offsetFactors(noffsets);
    for (size_t i = 0; i < noffsets; i++) {
      size_t k = 1;
      for (size_t j = 0; j < i; j++)
        k *= this->dimensions_[j];
      offsetFactors[i] = k;
    }

    // Calculate offset
    size_t offset = ind[0];
    for (size_t i = 1; i < ind.size(); i++) {
      offset += ind[i] * offsetFactors[i];
    }

    if (offset >= this->get_size()) {
      throw std::runtime_error("hoNDArray::operator() index out of bounds");
    }
    return this->data_[offset];
  }

  bool operator==(hoNDArray const& other) const {
    if (this->dimensions_ != other.dimensions_) {
      return false;
    }

    if (this->dimensions_.empty()) {
      return true;
    }

    auto nelements = this->get_size();

    for (size_t i = 0; i < nelements; i++) {
      if (this->data_[i] != other.data_[i]) {
        return false;
      }
    }
    return true;
  }

  bool operator!=(hoNDArray const& other) const {
    return !(*this == other);
  }

  T* begin() { return this->get_data(); }
  T* end() { return this->get_data() + this->get_size(); }
  T const* begin() const { return this->get_data(); }
  T const* end() const { return this->get_data() + this->get_size(); }

 protected:
  virtual void deallocate_memory() override {
    if (this->data_ != nullptr) {
      delete[] this->data_;
      this->data_ = 0x0;
    }
  }

  virtual void allocate_memory() override {
    deallocate_memory();

    if (!this->dimensions_.empty()) {
      auto nelements = this->get_size();

      if (nelements > 0) {
        this->data_ = new T[nelements];

        if (this->data_ == nullptr) {
          throw std::runtime_error("hoNDArray<>::allocate memory failed");
        }

        std::fill(this->data_, this->data_ + nelements, T{});
      }
    }
  }
};

}  // namespace external
