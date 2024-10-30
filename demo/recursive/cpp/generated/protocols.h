// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include "types.h"

namespace sketch {
enum class Version {
  Current
};
// Abstract writer for the MyProtocol protocol.
class MyProtocolWriterBase {
  public:
  // Ordinal 0.
  void WriteTree(sketch::BinaryTree const& value);

  // Ordinal 1.
  void WritePtree(std::unique_ptr<sketch::BinaryTree> const& value);

  // Ordinal 2.
  void WriteList(std::optional<sketch::LinkedList<std::string>> const& value);

  // Ordinal 3.
  // dirs: !stream
  //   items: Directory
  // Call this method for each element of the `cwd` stream, then call `EndCwd() when done.`
  void WriteCwd(sketch::DirectoryEntry const& value);

  // Ordinal 3.
  // dirs: !stream
  //   items: Directory
  // Call this method to write many values to the `cwd` stream, then call `EndCwd()` when done.
  void WriteCwd(std::vector<sketch::DirectoryEntry> const& values);

  // Marks the end of the `cwd` stream.
  void EndCwd();

  // Optionaly close this writer before destructing. Validates that all steps were completed.
  void Close();

  virtual ~MyProtocolWriterBase() = default;

  // Flushes all buffered data.
  virtual void Flush() {}

  protected:
  virtual void WriteTreeImpl(sketch::BinaryTree const& value) = 0;
  virtual void WritePtreeImpl(std::unique_ptr<sketch::BinaryTree> const& value) = 0;
  virtual void WriteListImpl(std::optional<sketch::LinkedList<std::string>> const& value) = 0;
  virtual void WriteCwdImpl(sketch::DirectoryEntry const& value) = 0;
  virtual void WriteCwdImpl(std::vector<sketch::DirectoryEntry> const& value);
  virtual void EndCwdImpl() = 0;
  virtual void CloseImpl() {}

  static std::string schema_;

  static std::vector<std::string> previous_schemas_;

  static std::string SchemaFromVersion(Version version);

  private:
  uint8_t state_ = 0;

  friend class MyProtocolReaderBase;
  friend class MyProtocolIndexedReaderBase;
};

// Abstract reader for the MyProtocol protocol.
class MyProtocolReaderBase {
  public:
  // Ordinal 0.
  void ReadTree(sketch::BinaryTree& value);

  // Ordinal 1.
  void ReadPtree(std::unique_ptr<sketch::BinaryTree>& value);

  // Ordinal 2.
  void ReadList(std::optional<sketch::LinkedList<std::string>>& value);

  // Ordinal 3.
  // dirs: !stream
  //   items: Directory
  [[nodiscard]] bool ReadCwd(sketch::DirectoryEntry& value);

  // Ordinal 3.
  // dirs: !stream
  //   items: Directory
  [[nodiscard]] bool ReadCwd(std::vector<sketch::DirectoryEntry>& values);

  // Optionaly close this writer before destructing. Validates that all steps were completely read.
  void Close();

  void CopyTo(MyProtocolWriterBase& writer, size_t cwd_buffer_size = 1);

  virtual ~MyProtocolReaderBase() = default;

  protected:
  virtual void ReadTreeImpl(sketch::BinaryTree& value) = 0;
  virtual void ReadPtreeImpl(std::unique_ptr<sketch::BinaryTree>& value) = 0;
  virtual void ReadListImpl(std::optional<sketch::LinkedList<std::string>>& value) = 0;
  virtual bool ReadCwdImpl(sketch::DirectoryEntry& value) = 0;
  virtual bool ReadCwdImpl(std::vector<sketch::DirectoryEntry>& values);
  virtual void CloseImpl() {}
  static std::string schema_;

  static std::vector<std::string> previous_schemas_;

  static Version VersionFromSchema(const std::string& schema);

  private:
  uint8_t state_ = 0;
};

// Abstract Indexed reader for the MyProtocol protocol.
class MyProtocolIndexedReaderBase {
  public:
  // Ordinal 0.
  void ReadTree(sketch::BinaryTree& value);

  // Ordinal 1.
  void ReadPtree(std::unique_ptr<sketch::BinaryTree>& value);

  // Ordinal 2.
  void ReadList(std::optional<sketch::LinkedList<std::string>>& value);

  // Ordinal 3.
  // dirs: !stream
  //   items: Directory
  [[nodiscard]] bool ReadCwd(sketch::DirectoryEntry& value, size_t idx=0);

  // Ordinal 3.
  // dirs: !stream
  //   items: Directory
  [[nodiscard]] bool ReadCwd(std::vector<sketch::DirectoryEntry>& values, size_t idx=0);

  [[nodiscard]] size_t CountCwd();

  // Optionaly close this writer before destructing
  void Close();

  virtual ~MyProtocolIndexedReaderBase() = default;

  protected:
  virtual void ReadTreeImpl(sketch::BinaryTree& value) = 0;
  virtual void ReadPtreeImpl(std::unique_ptr<sketch::BinaryTree>& value) = 0;
  virtual void ReadListImpl(std::optional<sketch::LinkedList<std::string>>& value) = 0;
  virtual bool ReadCwdImpl(sketch::DirectoryEntry& value, size_t idx) = 0;
  virtual bool ReadCwdImpl(std::vector<sketch::DirectoryEntry>& values, size_t idx) = 0;
  virtual size_t CountCwdImpl() = 0;
  virtual void CloseImpl() {}
  static std::string schema_;

  static std::vector<std::string> previous_schemas_;

  static Version VersionFromSchema(const std::string& schema);

  private:
  uint8_t state_ = 0;
};
} // namespace sketch
