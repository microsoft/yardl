// This file was generated by the "yardl" tool. DO NOT EDIT.

#include "protocols.h"

#include "../yardl/detail/hdf5/io.h"
#include "../yardl/detail/hdf5/ddl.h"
#include "../yardl/detail/hdf5/inner_types.h"

namespace {
template <typename TInner0, typename TOuter0, typename TInner1, typename TOuter1>
class InnerUnion2 {
  public:
  InnerUnion2() : type_index_(-1) {} 
  InnerUnion2(std::variant<TOuter0, TOuter1> const& v) : type_index_(static_cast<int8_t>(v.index())) {
    Init(v);
  }

  InnerUnion2(std::variant<std::monostate, TOuter0, TOuter1> const& v) : type_index_(static_cast<int8_t>(v.index()) - 1) {
    Init(v);
  }

  InnerUnion2(InnerUnion2 const& v) = delete;

  InnerUnion2 operator=(InnerUnion2 const&) = delete;

  ~InnerUnion2() {
    switch (type_index_) {
    case 0:
      value0_.~TInner0();
      break;
    case 1:
      value1_.~TInner1();
      break;
    }
  }

  void ToOuter(std::variant<TOuter0, TOuter1>& o) const {
    ToOuterImpl(o);
  }

  void ToOuter(std::variant<std::monostate, TOuter0, TOuter1>& o) const {
    ToOuterImpl(o);
  }

  int8_t type_index_;
  union {
    char empty0_[sizeof(TInner0)]{};
    TInner0 value0_;
  };
  union {
    char empty1_[sizeof(TInner1)]{};
    TInner1 value1_;
  };

  private:
  template <typename T>
  void Init(T const& v) {
    constexpr size_t offset = GetOuterVariantOffset<std::remove_const_t<std::remove_reference_t<decltype(v)>>>();
    switch (type_index_) {
    case 0:
      new (&value0_) TInner0(std::get<0 + offset>(v));
      return;
    case 1:
      new (&value1_) TInner1(std::get<1 + offset>(v));
      return;
    }
  }

  template <typename TVariant>
  void ToOuterImpl(TVariant& o) const {
    constexpr size_t offset = GetOuterVariantOffset<TVariant>();
    switch (type_index_) {
    case -1:
      if constexpr (offset == 1) {
        o.template emplace<0>(std::monostate{});
        return;
      }
    case 0:
      o.template emplace<0 + offset>();
      yardl::hdf5::ToOuter(value0_, std::get<0 + offset>(o));
      return;
    case 1:
      o.template emplace<1 + offset>();
      yardl::hdf5::ToOuter(value1_, std::get<1 + offset>(o));
      return;
    }
    throw std::runtime_error("unrecognized type variant type index " + std::to_string(type_index_));
  }

  template <typename TVariant>
  static constexpr size_t GetOuterVariantOffset() {
    constexpr bool has_monostate = std::is_same_v<std::monostate, std::variant_alternative_t<0, TVariant>>;
    if constexpr (has_monostate) {
      return 1;
    }
      return 0;
  }
};

template <typename TInner0, typename TOuter0, typename TInner1, typename TOuter1>
H5::CompType InnerUnion2Ddl(bool nullable, H5::DataType const& t0, std::string const& tag0, H5::DataType const& t1, std::string const& tag1) {
  using UnionType = ::InnerUnion2<TInner0, TOuter0, TInner1, TOuter1>;
  H5::CompType rtn(sizeof(UnionType));
  rtn.insertMember("$type", HOFFSET(UnionType, type_index_), yardl::hdf5::UnionTypeEnumDdl(nullable, tag0, tag1));
  rtn.insertMember(tag0, HOFFSET(UnionType, value0_), t0);
  rtn.insertMember(tag1, HOFFSET(UnionType, value1_), t1);
  return rtn;
}
}

namespace sketch::hdf5 {
namespace {
template <typename _T_Inner, typename T>
struct _Inner_LinkedList {
  _Inner_LinkedList() {} 
  _Inner_LinkedList(sketch::LinkedList<T> const& o) 
      : value(o.value),
      next(o.next) {
  }

  void ToOuter (sketch::LinkedList<T>& o) const {
    yardl::hdf5::ToOuter(value, o.value);
    yardl::hdf5::ToOuter(next, o.next);
  }

  _T_Inner value;
  sketch::hdf5::_Inner_LinkedList<_T_Inner, T> next;
};

struct _Inner_File {
  _Inner_File() {} 
  _Inner_File(sketch::File const& o) 
      : name(o.name),
      data(o.data) {
  }

  void ToOuter (sketch::File& o) const {
    yardl::hdf5::ToOuter(name, o.name);
    yardl::hdf5::ToOuter(data, o.data);
  }

  yardl::hdf5::InnerVlenString name;
  yardl::hdf5::InnerVlen<uint8_t, uint8_t> data;
};

struct _Inner_Directory {
  _Inner_Directory() {} 
  _Inner_Directory(sketch::Directory const& o) 
      : name(o.name),
      entries(o.entries) {
  }

  void ToOuter (sketch::Directory& o) const {
    yardl::hdf5::ToOuter(name, o.name);
    yardl::hdf5::ToOuter(entries, o.entries);
  }

  yardl::hdf5::InnerVlenString name;
  yardl::hdf5::InnerVlen<::InnerUnion2<sketch::hdf5::_Inner_File, sketch::File, sketch::hdf5::_Inner_Directory, std::unique_ptr<sketch::Directory>>, sketch::DirectoryEntry> entries;
};

[[maybe_unused]] H5::CompType GetBinaryTreeHdf5Ddl() {
  using RecordType = sketch::BinaryTree;
  H5::CompType t(sizeof(RecordType));
  t.insertMember("value", HOFFSET(RecordType, value), H5::PredType::NATIVE_INT32);
  t.insertMember("left", HOFFSET(RecordType, left), sketch::hdf5::GetBinaryTreeHdf5Ddl());
  t.insertMember("right", HOFFSET(RecordType, right), sketch::hdf5::GetBinaryTreeHdf5Ddl());
  return t;
}

template <typename _T_Inner, typename T>
[[maybe_unused]] H5::CompType GetLinkedListHdf5Ddl(H5::DataType const& T_type) {
  using RecordType = sketch::hdf5::_Inner_LinkedList<_T_Inner, T>;
  H5::CompType t(sizeof(RecordType));
  t.insertMember("value", HOFFSET(RecordType, value), T_type);
  t.insertMember("next", HOFFSET(RecordType, next), sketch::hdf5::GetLinkedListHdf5Ddl<_T_Inner, T>(T_type));
  return t;
}

[[maybe_unused]] H5::CompType GetFileHdf5Ddl() {
  using RecordType = sketch::hdf5::_Inner_File;
  H5::CompType t(sizeof(RecordType));
  t.insertMember("name", HOFFSET(RecordType, name), yardl::hdf5::InnerVlenStringDdl());
  t.insertMember("data", HOFFSET(RecordType, data), yardl::hdf5::InnerVlenDdl(H5::PredType::NATIVE_UINT8));
  return t;
}

[[maybe_unused]] H5::CompType GetDirectoryHdf5Ddl() {
  using RecordType = sketch::hdf5::_Inner_Directory;
  H5::CompType t(sizeof(RecordType));
  t.insertMember("name", HOFFSET(RecordType, name), yardl::hdf5::InnerVlenStringDdl());
  t.insertMember("entries", HOFFSET(RecordType, entries), yardl::hdf5::InnerVlenDdl(::InnerUnion2Ddl<sketch::hdf5::_Inner_File, sketch::File, sketch::hdf5::_Inner_Directory, std::unique_ptr<sketch::Directory>>(false, sketch::hdf5::GetFileHdf5Ddl(), "File", sketch::hdf5::GetDirectoryHdf5Ddl(), "Directory")));
  return t;
}

} // namespace 

MyProtocolWriter::MyProtocolWriter(std::string path)
    : yardl::hdf5::Hdf5Writer::Hdf5Writer(path, "MyProtocol", schema_) {
}

void MyProtocolWriter::WriteTreeImpl(sketch::BinaryTree const& value) {
  yardl::hdf5::WriteScalarDataset<sketch::BinaryTree, sketch::BinaryTree>(group_, "tree", sketch::hdf5::GetBinaryTreeHdf5Ddl(), value);
}

void MyProtocolWriter::WritePtreeImpl(std::unique_ptr<sketch::BinaryTree> const& value) {
  yardl::hdf5::WriteScalarDataset<sketch::BinaryTree, std::unique_ptr<sketch::BinaryTree>>(group_, "ptree", sketch::hdf5::GetBinaryTreeHdf5Ddl(), value);
}

void MyProtocolWriter::WriteTreesImpl(sketch::BinaryTree const& value) {
  if (!trees_dataset_state_) {
    trees_dataset_state_ = std::make_unique<yardl::hdf5::DatasetWriter>(group_, "trees", sketch::hdf5::GetBinaryTreeHdf5Ddl(), 0);
  }

  trees_dataset_state_->Append<sketch::BinaryTree, sketch::BinaryTree>(value);
}

void MyProtocolWriter::WriteTreesImpl(std::vector<sketch::BinaryTree> const& values) {
  if (!trees_dataset_state_) {
    trees_dataset_state_ = std::make_unique<yardl::hdf5::DatasetWriter>(group_, "trees", sketch::hdf5::GetBinaryTreeHdf5Ddl(), 0);
  }

  trees_dataset_state_->AppendBatch<sketch::BinaryTree, sketch::BinaryTree>(values);
}

void MyProtocolWriter::EndTreesImpl() {
  if (!trees_dataset_state_) {
    trees_dataset_state_ = std::make_unique<yardl::hdf5::DatasetWriter>(group_, "trees", sketch::hdf5::GetBinaryTreeHdf5Ddl(), 0);
  }

  trees_dataset_state_.reset();
}

void MyProtocolWriter::WritePtreesImpl(std::unique_ptr<sketch::BinaryTree> const& value) {
  if (!ptrees_dataset_state_) {
    ptrees_dataset_state_ = std::make_unique<yardl::hdf5::DatasetWriter>(group_, "ptrees", sketch::hdf5::GetBinaryTreeHdf5Ddl(), 0);
  }

  ptrees_dataset_state_->Append<sketch::BinaryTree, std::unique_ptr<sketch::BinaryTree>>(value);
}

void MyProtocolWriter::WritePtreesImpl(std::vector<std::unique_ptr<sketch::BinaryTree>> const& values) {
  if (!ptrees_dataset_state_) {
    ptrees_dataset_state_ = std::make_unique<yardl::hdf5::DatasetWriter>(group_, "ptrees", sketch::hdf5::GetBinaryTreeHdf5Ddl(), 0);
  }

  ptrees_dataset_state_->AppendBatch<sketch::BinaryTree, std::unique_ptr<sketch::BinaryTree>>(values);
}

void MyProtocolWriter::EndPtreesImpl() {
  if (!ptrees_dataset_state_) {
    ptrees_dataset_state_ = std::make_unique<yardl::hdf5::DatasetWriter>(group_, "ptrees", sketch::hdf5::GetBinaryTreeHdf5Ddl(), 0);
  }

  ptrees_dataset_state_.reset();
}

void MyProtocolWriter::WriteListImpl(sketch::LinkedList<int32_t> const& value) {
  yardl::hdf5::WriteScalarDataset<sketch::hdf5::_Inner_LinkedList<int32_t, int32_t>, sketch::LinkedList<int32_t>>(group_, "list", sketch::hdf5::GetLinkedListHdf5Ddl<int32_t, int32_t>(H5::PredType::NATIVE_INT32), value);
}

MyProtocolReader::MyProtocolReader(std::string path)
    : yardl::hdf5::Hdf5Reader::Hdf5Reader(path, "MyProtocol", schema_) {
}

void MyProtocolReader::ReadTreeImpl(sketch::BinaryTree& value) {
  yardl::hdf5::ReadScalarDataset<sketch::BinaryTree, sketch::BinaryTree>(group_, "tree", sketch::hdf5::GetBinaryTreeHdf5Ddl(), value);
}

void MyProtocolReader::ReadPtreeImpl(std::unique_ptr<sketch::BinaryTree>& value) {
  yardl::hdf5::ReadScalarDataset<sketch::BinaryTree, std::unique_ptr<sketch::BinaryTree>>(group_, "ptree", sketch::hdf5::GetBinaryTreeHdf5Ddl(), value);
}

bool MyProtocolReader::ReadTreesImpl(sketch::BinaryTree& value) {
  if (!trees_dataset_state_) {
    trees_dataset_state_ = std::make_unique<yardl::hdf5::DatasetReader>(group_, "trees", sketch::hdf5::GetBinaryTreeHdf5Ddl(), 0);
  }

  bool has_value = trees_dataset_state_->Read<sketch::BinaryTree, sketch::BinaryTree>(value);
  if (!has_value) {
    trees_dataset_state_.reset();
  }

  return has_value;
}

bool MyProtocolReader::ReadTreesImpl(std::vector<sketch::BinaryTree>& values) {
  if (!trees_dataset_state_) {
    trees_dataset_state_ = std::make_unique<yardl::hdf5::DatasetReader>(group_, "trees", sketch::hdf5::GetBinaryTreeHdf5Ddl());
  }

  bool has_more = trees_dataset_state_->ReadBatch<sketch::BinaryTree, sketch::BinaryTree>(values);
  if (!has_more) {
    trees_dataset_state_.reset();
  }

  return has_more;
}

bool MyProtocolReader::ReadPtreesImpl(std::unique_ptr<sketch::BinaryTree>& value) {
  if (!ptrees_dataset_state_) {
    ptrees_dataset_state_ = std::make_unique<yardl::hdf5::DatasetReader>(group_, "ptrees", sketch::hdf5::GetBinaryTreeHdf5Ddl(), 0);
  }

  bool has_value = ptrees_dataset_state_->Read<sketch::BinaryTree, std::unique_ptr<sketch::BinaryTree>>(value);
  if (!has_value) {
    ptrees_dataset_state_.reset();
  }

  return has_value;
}

bool MyProtocolReader::ReadPtreesImpl(std::vector<std::unique_ptr<sketch::BinaryTree>>& values) {
  if (!ptrees_dataset_state_) {
    ptrees_dataset_state_ = std::make_unique<yardl::hdf5::DatasetReader>(group_, "ptrees", sketch::hdf5::GetBinaryTreeHdf5Ddl());
  }

  bool has_more = ptrees_dataset_state_->ReadBatch<sketch::BinaryTree, std::unique_ptr<sketch::BinaryTree>>(values);
  if (!has_more) {
    ptrees_dataset_state_.reset();
  }

  return has_more;
}

void MyProtocolReader::ReadListImpl(sketch::LinkedList<int32_t>& value) {
  yardl::hdf5::ReadScalarDataset<sketch::hdf5::_Inner_LinkedList<int32_t, int32_t>, sketch::LinkedList<int32_t>>(group_, "list", sketch::hdf5::GetLinkedListHdf5Ddl<int32_t, int32_t>(H5::PredType::NATIVE_INT32), value);
}

} // namespace sketch::hdf5

