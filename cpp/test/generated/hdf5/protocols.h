// This file was generated by the "yardl" tool. DO NOT EDIT.

#pragma once
#include <array>
#include <complex>
#include <optional>
#include <variant>
#include <vector>

#include "../protocols.h"
#include "../types.h"
#include "../yardl/detail/hdf5/io.h"

namespace test_model::hdf5 {
// HDF5 writer for the BenchmarkFloat256x256 protocol.
class BenchmarkFloat256x256Writer : public test_model::BenchmarkFloat256x256WriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  BenchmarkFloat256x256Writer(std::string path);

  protected:
  void WriteFloat256x256Impl(yardl::FixedNDArray<float, 256, 256> const& value) override;

  void WriteFloat256x256Impl(std::vector<yardl::FixedNDArray<float, 256, 256>> const& values) override;

  void EndFloat256x256Impl() override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetWriter> float256x256_dataset_state_;
};

// HDF5 reader for the BenchmarkFloat256x256 protocol.
class BenchmarkFloat256x256Reader : public test_model::BenchmarkFloat256x256ReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  BenchmarkFloat256x256Reader(std::string path);

  bool ReadFloat256x256Impl(yardl::FixedNDArray<float, 256, 256>& value) override;

  bool ReadFloat256x256Impl(std::vector<yardl::FixedNDArray<float, 256, 256>>& values) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetReader> float256x256_dataset_state_;
};
// HDF5 writer for the BenchmarkFloatVlen protocol.
class BenchmarkFloatVlenWriter : public test_model::BenchmarkFloatVlenWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  BenchmarkFloatVlenWriter(std::string path);

  protected:
  void WriteFloatArrayImpl(yardl::NDArray<float, 2> const& value) override;

  void WriteFloatArrayImpl(std::vector<yardl::NDArray<float, 2>> const& values) override;

  void EndFloatArrayImpl() override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetWriter> floatArray_dataset_state_;
};

// HDF5 reader for the BenchmarkFloatVlen protocol.
class BenchmarkFloatVlenReader : public test_model::BenchmarkFloatVlenReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  BenchmarkFloatVlenReader(std::string path);

  bool ReadFloatArrayImpl(yardl::NDArray<float, 2>& value) override;

  bool ReadFloatArrayImpl(std::vector<yardl::NDArray<float, 2>>& values) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetReader> floatArray_dataset_state_;
};
// HDF5 writer for the BenchmarkSmallRecord protocol.
class BenchmarkSmallRecordWriter : public test_model::BenchmarkSmallRecordWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  BenchmarkSmallRecordWriter(std::string path);

  protected:
  void WriteSmallRecordImpl(test_model::SmallBenchmarkRecord const& value) override;

  void WriteSmallRecordImpl(std::vector<test_model::SmallBenchmarkRecord> const& values) override;

  void EndSmallRecordImpl() override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetWriter> smallRecord_dataset_state_;
};

// HDF5 reader for the BenchmarkSmallRecord protocol.
class BenchmarkSmallRecordReader : public test_model::BenchmarkSmallRecordReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  BenchmarkSmallRecordReader(std::string path);

  bool ReadSmallRecordImpl(test_model::SmallBenchmarkRecord& value) override;

  bool ReadSmallRecordImpl(std::vector<test_model::SmallBenchmarkRecord>& values) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetReader> smallRecord_dataset_state_;
};
// HDF5 writer for the BenchmarkSmallRecordWithOptionals protocol.
class BenchmarkSmallRecordWithOptionalsWriter : public test_model::BenchmarkSmallRecordWithOptionalsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  BenchmarkSmallRecordWithOptionalsWriter(std::string path);

  protected:
  void WriteSmallRecordImpl(test_model::SimpleEncodingCounters const& value) override;

  void WriteSmallRecordImpl(std::vector<test_model::SimpleEncodingCounters> const& values) override;

  void EndSmallRecordImpl() override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetWriter> smallRecord_dataset_state_;
};

// HDF5 reader for the BenchmarkSmallRecordWithOptionals protocol.
class BenchmarkSmallRecordWithOptionalsReader : public test_model::BenchmarkSmallRecordWithOptionalsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  BenchmarkSmallRecordWithOptionalsReader(std::string path);

  bool ReadSmallRecordImpl(test_model::SimpleEncodingCounters& value) override;

  bool ReadSmallRecordImpl(std::vector<test_model::SimpleEncodingCounters>& values) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetReader> smallRecord_dataset_state_;
};
// HDF5 writer for the BenchmarkSimpleMrd protocol.
class BenchmarkSimpleMrdWriter : public test_model::BenchmarkSimpleMrdWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  BenchmarkSimpleMrdWriter(std::string path);

  protected:
  void WriteDataImpl(std::variant<test_model::SimpleAcquisition, test_model::Image<float>> const& value) override;

  void EndDataImpl() override;

  public:
  void Flush() override;

  private:
  std::unique_ptr<yardl::hdf5::UnionDatasetWriter<2>> data_dataset_state_;
};

// HDF5 reader for the BenchmarkSimpleMrd protocol.
class BenchmarkSimpleMrdReader : public test_model::BenchmarkSimpleMrdReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  BenchmarkSimpleMrdReader(std::string path);

  bool ReadDataImpl(std::variant<test_model::SimpleAcquisition, test_model::Image<float>>& value) override;

  private:
  std::unique_ptr<yardl::hdf5::UnionDatasetReader<2>> data_dataset_state_;
};
// HDF5 writer for the Scalars protocol.
class ScalarsWriter : public test_model::ScalarsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  ScalarsWriter(std::string path);

  protected:
  void WriteInt32Impl(int32_t const& value) override;

  void WriteRecordImpl(test_model::RecordWithPrimitives const& value) override;

  private:
};

// HDF5 reader for the Scalars protocol.
class ScalarsReader : public test_model::ScalarsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  ScalarsReader(std::string path);

  void ReadInt32Impl(int32_t& value) override;

  void ReadRecordImpl(test_model::RecordWithPrimitives& value) override;

  private:
};
// HDF5 writer for the ScalarOptionals protocol.
class ScalarOptionalsWriter : public test_model::ScalarOptionalsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  ScalarOptionalsWriter(std::string path);

  protected:
  void WriteOptionalIntImpl(std::optional<int32_t> const& value) override;

  void WriteOptionalRecordImpl(std::optional<test_model::SimpleRecord> const& value) override;

  void WriteRecordWithOptionalFieldsImpl(test_model::RecordWithOptionalFields const& value) override;

  void WriteOptionalRecordWithOptionalFieldsImpl(std::optional<test_model::RecordWithOptionalFields> const& value) override;

  private:
};

// HDF5 reader for the ScalarOptionals protocol.
class ScalarOptionalsReader : public test_model::ScalarOptionalsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  ScalarOptionalsReader(std::string path);

  void ReadOptionalIntImpl(std::optional<int32_t>& value) override;

  void ReadOptionalRecordImpl(std::optional<test_model::SimpleRecord>& value) override;

  void ReadRecordWithOptionalFieldsImpl(test_model::RecordWithOptionalFields& value) override;

  void ReadOptionalRecordWithOptionalFieldsImpl(std::optional<test_model::RecordWithOptionalFields>& value) override;

  private:
};
// HDF5 writer for the NestedRecords protocol.
class NestedRecordsWriter : public test_model::NestedRecordsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  NestedRecordsWriter(std::string path);

  protected:
  void WriteTupleWithRecordsImpl(test_model::TupleWithRecords const& value) override;

  private:
};

// HDF5 reader for the NestedRecords protocol.
class NestedRecordsReader : public test_model::NestedRecordsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  NestedRecordsReader(std::string path);

  void ReadTupleWithRecordsImpl(test_model::TupleWithRecords& value) override;

  private:
};
// HDF5 writer for the Vlens protocol.
class VlensWriter : public test_model::VlensWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  VlensWriter(std::string path);

  protected:
  void WriteIntVectorImpl(std::vector<int32_t> const& value) override;

  void WriteComplexVectorImpl(std::vector<std::complex<float>> const& value) override;

  void WriteRecordWithVlensImpl(test_model::RecordWithVlens const& value) override;

  void WriteVlenOfRecordWithVlensImpl(std::vector<test_model::RecordWithVlens> const& value) override;

  private:
};

// HDF5 reader for the Vlens protocol.
class VlensReader : public test_model::VlensReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  VlensReader(std::string path);

  void ReadIntVectorImpl(std::vector<int32_t>& value) override;

  void ReadComplexVectorImpl(std::vector<std::complex<float>>& value) override;

  void ReadRecordWithVlensImpl(test_model::RecordWithVlens& value) override;

  void ReadVlenOfRecordWithVlensImpl(std::vector<test_model::RecordWithVlens>& value) override;

  private:
};
// HDF5 writer for the Strings protocol.
class StringsWriter : public test_model::StringsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  StringsWriter(std::string path);

  protected:
  void WriteSingleStringImpl(std::string const& value) override;

  void WriteRecWithStringImpl(test_model::RecordWithStrings const& value) override;

  private:
};

// HDF5 reader for the Strings protocol.
class StringsReader : public test_model::StringsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  StringsReader(std::string path);

  void ReadSingleStringImpl(std::string& value) override;

  void ReadRecWithStringImpl(test_model::RecordWithStrings& value) override;

  private:
};
// HDF5 writer for the OptionalVectors protocol.
class OptionalVectorsWriter : public test_model::OptionalVectorsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  OptionalVectorsWriter(std::string path);

  protected:
  void WriteRecordWithOptionalVectorImpl(test_model::RecordWithOptionalVector const& value) override;

  private:
};

// HDF5 reader for the OptionalVectors protocol.
class OptionalVectorsReader : public test_model::OptionalVectorsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  OptionalVectorsReader(std::string path);

  void ReadRecordWithOptionalVectorImpl(test_model::RecordWithOptionalVector& value) override;

  private:
};
// HDF5 writer for the FixedVectors protocol.
class FixedVectorsWriter : public test_model::FixedVectorsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  FixedVectorsWriter(std::string path);

  protected:
  void WriteFixedIntVectorImpl(std::array<int32_t, 5> const& value) override;

  void WriteFixedSimpleRecordVectorImpl(std::array<test_model::SimpleRecord, 3> const& value) override;

  void WriteFixedRecordWithVlensVectorImpl(std::array<test_model::RecordWithVlens, 2> const& value) override;

  void WriteRecordWithFixedVectorsImpl(test_model::RecordWithFixedVectors const& value) override;

  private:
};

// HDF5 reader for the FixedVectors protocol.
class FixedVectorsReader : public test_model::FixedVectorsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  FixedVectorsReader(std::string path);

  void ReadFixedIntVectorImpl(std::array<int32_t, 5>& value) override;

  void ReadFixedSimpleRecordVectorImpl(std::array<test_model::SimpleRecord, 3>& value) override;

  void ReadFixedRecordWithVlensVectorImpl(std::array<test_model::RecordWithVlens, 2>& value) override;

  void ReadRecordWithFixedVectorsImpl(test_model::RecordWithFixedVectors& value) override;

  private:
};
// HDF5 writer for the Streams protocol.
class StreamsWriter : public test_model::StreamsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  StreamsWriter(std::string path);

  protected:
  void WriteIntDataImpl(int32_t const& value) override;

  void WriteIntDataImpl(std::vector<int32_t> const& values) override;

  void EndIntDataImpl() override;

  void WriteOptionalIntDataImpl(std::optional<int32_t> const& value) override;

  void WriteOptionalIntDataImpl(std::vector<std::optional<int32_t>> const& values) override;

  void EndOptionalIntDataImpl() override;

  void WriteRecordWithOptionalVectorDataImpl(test_model::RecordWithOptionalVector const& value) override;

  void WriteRecordWithOptionalVectorDataImpl(std::vector<test_model::RecordWithOptionalVector> const& values) override;

  void EndRecordWithOptionalVectorDataImpl() override;

  void WriteFixedVectorImpl(std::array<int32_t, 3> const& value) override;

  void WriteFixedVectorImpl(std::vector<std::array<int32_t, 3>> const& values) override;

  void EndFixedVectorImpl() override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetWriter> intData_dataset_state_;
  std::unique_ptr<yardl::hdf5::DatasetWriter> optionalIntData_dataset_state_;
  std::unique_ptr<yardl::hdf5::DatasetWriter> recordWithOptionalVectorData_dataset_state_;
  std::unique_ptr<yardl::hdf5::DatasetWriter> fixedVector_dataset_state_;
};

// HDF5 reader for the Streams protocol.
class StreamsReader : public test_model::StreamsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  StreamsReader(std::string path);

  bool ReadIntDataImpl(int32_t& value) override;

  bool ReadIntDataImpl(std::vector<int32_t>& values) override;

  bool ReadOptionalIntDataImpl(std::optional<int32_t>& value) override;

  bool ReadOptionalIntDataImpl(std::vector<std::optional<int32_t>>& values) override;

  bool ReadRecordWithOptionalVectorDataImpl(test_model::RecordWithOptionalVector& value) override;

  bool ReadRecordWithOptionalVectorDataImpl(std::vector<test_model::RecordWithOptionalVector>& values) override;

  bool ReadFixedVectorImpl(std::array<int32_t, 3>& value) override;

  bool ReadFixedVectorImpl(std::vector<std::array<int32_t, 3>>& values) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetReader> intData_dataset_state_;
  std::unique_ptr<yardl::hdf5::DatasetReader> optionalIntData_dataset_state_;
  std::unique_ptr<yardl::hdf5::DatasetReader> recordWithOptionalVectorData_dataset_state_;
  std::unique_ptr<yardl::hdf5::DatasetReader> fixedVector_dataset_state_;
};
// HDF5 writer for the FixedArrays protocol.
class FixedArraysWriter : public test_model::FixedArraysWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  FixedArraysWriter(std::string path);

  protected:
  void WriteIntsImpl(yardl::FixedNDArray<int32_t, 2, 3> const& value) override;

  void WriteFixedSimpleRecordArrayImpl(yardl::FixedNDArray<test_model::SimpleRecord, 3, 2> const& value) override;

  void WriteFixedRecordWithVlensArrayImpl(yardl::FixedNDArray<test_model::RecordWithVlens, 2, 2> const& value) override;

  void WriteRecordWithFixedArraysImpl(test_model::RecordWithFixedArrays const& value) override;

  void WriteNamedArrayImpl(test_model::NamedFixedNDArray const& value) override;

  private:
};

// HDF5 reader for the FixedArrays protocol.
class FixedArraysReader : public test_model::FixedArraysReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  FixedArraysReader(std::string path);

  void ReadIntsImpl(yardl::FixedNDArray<int32_t, 2, 3>& value) override;

  void ReadFixedSimpleRecordArrayImpl(yardl::FixedNDArray<test_model::SimpleRecord, 3, 2>& value) override;

  void ReadFixedRecordWithVlensArrayImpl(yardl::FixedNDArray<test_model::RecordWithVlens, 2, 2>& value) override;

  void ReadRecordWithFixedArraysImpl(test_model::RecordWithFixedArrays& value) override;

  void ReadNamedArrayImpl(test_model::NamedFixedNDArray& value) override;

  private:
};
// HDF5 writer for the Subarrays protocol.
class SubarraysWriter : public test_model::SubarraysWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  SubarraysWriter(std::string path);

  protected:
  void WriteDynamicWithFixedIntSubarrayImpl(yardl::DynamicNDArray<yardl::FixedNDArray<int32_t, 3>> const& value) override;

  void WriteDynamicWithFixedFloatSubarrayImpl(yardl::DynamicNDArray<yardl::FixedNDArray<float, 3>> const& value) override;

  void WriteKnownDimCountWithFixedIntSubarrayImpl(yardl::NDArray<yardl::FixedNDArray<int32_t, 3>, 1> const& value) override;

  void WriteKnownDimCountWithFixedFloatSubarrayImpl(yardl::NDArray<yardl::FixedNDArray<float, 3>, 1> const& value) override;

  void WriteFixedWithFixedIntSubarrayImpl(yardl::FixedNDArray<yardl::FixedNDArray<int32_t, 3>, 2> const& value) override;

  void WriteFixedWithFixedFloatSubarrayImpl(yardl::FixedNDArray<yardl::FixedNDArray<float, 3>, 2> const& value) override;

  void WriteNestedSubarrayImpl(yardl::DynamicNDArray<yardl::FixedNDArray<yardl::FixedNDArray<int32_t, 3>, 2>> const& value) override;

  void WriteDynamicWithFixedVectorSubarrayImpl(yardl::DynamicNDArray<std::array<int32_t, 3>> const& value) override;

  private:
};

// HDF5 reader for the Subarrays protocol.
class SubarraysReader : public test_model::SubarraysReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  SubarraysReader(std::string path);

  void ReadDynamicWithFixedIntSubarrayImpl(yardl::DynamicNDArray<yardl::FixedNDArray<int32_t, 3>>& value) override;

  void ReadDynamicWithFixedFloatSubarrayImpl(yardl::DynamicNDArray<yardl::FixedNDArray<float, 3>>& value) override;

  void ReadKnownDimCountWithFixedIntSubarrayImpl(yardl::NDArray<yardl::FixedNDArray<int32_t, 3>, 1>& value) override;

  void ReadKnownDimCountWithFixedFloatSubarrayImpl(yardl::NDArray<yardl::FixedNDArray<float, 3>, 1>& value) override;

  void ReadFixedWithFixedIntSubarrayImpl(yardl::FixedNDArray<yardl::FixedNDArray<int32_t, 3>, 2>& value) override;

  void ReadFixedWithFixedFloatSubarrayImpl(yardl::FixedNDArray<yardl::FixedNDArray<float, 3>, 2>& value) override;

  void ReadNestedSubarrayImpl(yardl::DynamicNDArray<yardl::FixedNDArray<yardl::FixedNDArray<int32_t, 3>, 2>>& value) override;

  void ReadDynamicWithFixedVectorSubarrayImpl(yardl::DynamicNDArray<std::array<int32_t, 3>>& value) override;

  private:
};
// HDF5 writer for the SubarraysInRecords protocol.
class SubarraysInRecordsWriter : public test_model::SubarraysInRecordsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  SubarraysInRecordsWriter(std::string path);

  protected:
  void WriteWithFixedSubarraysImpl(yardl::DynamicNDArray<test_model::RecordWithFixedCollections> const& value) override;

  void WriteWithVlenSubarraysImpl(yardl::DynamicNDArray<test_model::RecordWithVlenCollections> const& value) override;

  private:
};

// HDF5 reader for the SubarraysInRecords protocol.
class SubarraysInRecordsReader : public test_model::SubarraysInRecordsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  SubarraysInRecordsReader(std::string path);

  void ReadWithFixedSubarraysImpl(yardl::DynamicNDArray<test_model::RecordWithFixedCollections>& value) override;

  void ReadWithVlenSubarraysImpl(yardl::DynamicNDArray<test_model::RecordWithVlenCollections>& value) override;

  private:
};
// HDF5 writer for the NDArrays protocol.
class NDArraysWriter : public test_model::NDArraysWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  NDArraysWriter(std::string path);

  protected:
  void WriteIntsImpl(yardl::NDArray<int32_t, 2> const& value) override;

  void WriteSimpleRecordArrayImpl(yardl::NDArray<test_model::SimpleRecord, 2> const& value) override;

  void WriteRecordWithVlensArrayImpl(yardl::NDArray<test_model::RecordWithVlens, 2> const& value) override;

  void WriteRecordWithNDArraysImpl(test_model::RecordWithNDArrays const& value) override;

  void WriteNamedArrayImpl(test_model::NamedNDArray const& value) override;

  private:
};

// HDF5 reader for the NDArrays protocol.
class NDArraysReader : public test_model::NDArraysReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  NDArraysReader(std::string path);

  void ReadIntsImpl(yardl::NDArray<int32_t, 2>& value) override;

  void ReadSimpleRecordArrayImpl(yardl::NDArray<test_model::SimpleRecord, 2>& value) override;

  void ReadRecordWithVlensArrayImpl(yardl::NDArray<test_model::RecordWithVlens, 2>& value) override;

  void ReadRecordWithNDArraysImpl(test_model::RecordWithNDArrays& value) override;

  void ReadNamedArrayImpl(test_model::NamedNDArray& value) override;

  private:
};
// HDF5 writer for the NDArraysSingleDimension protocol.
class NDArraysSingleDimensionWriter : public test_model::NDArraysSingleDimensionWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  NDArraysSingleDimensionWriter(std::string path);

  protected:
  void WriteIntsImpl(yardl::NDArray<int32_t, 1> const& value) override;

  void WriteSimpleRecordArrayImpl(yardl::NDArray<test_model::SimpleRecord, 1> const& value) override;

  void WriteRecordWithVlensArrayImpl(yardl::NDArray<test_model::RecordWithVlens, 1> const& value) override;

  void WriteRecordWithNDArraysImpl(test_model::RecordWithNDArraysSingleDimension const& value) override;

  private:
};

// HDF5 reader for the NDArraysSingleDimension protocol.
class NDArraysSingleDimensionReader : public test_model::NDArraysSingleDimensionReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  NDArraysSingleDimensionReader(std::string path);

  void ReadIntsImpl(yardl::NDArray<int32_t, 1>& value) override;

  void ReadSimpleRecordArrayImpl(yardl::NDArray<test_model::SimpleRecord, 1>& value) override;

  void ReadRecordWithVlensArrayImpl(yardl::NDArray<test_model::RecordWithVlens, 1>& value) override;

  void ReadRecordWithNDArraysImpl(test_model::RecordWithNDArraysSingleDimension& value) override;

  private:
};
// HDF5 writer for the DynamicNDArrays protocol.
class DynamicNDArraysWriter : public test_model::DynamicNDArraysWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  DynamicNDArraysWriter(std::string path);

  protected:
  void WriteIntsImpl(yardl::DynamicNDArray<int32_t> const& value) override;

  void WriteSimpleRecordArrayImpl(yardl::DynamicNDArray<test_model::SimpleRecord> const& value) override;

  void WriteRecordWithVlensArrayImpl(yardl::DynamicNDArray<test_model::RecordWithVlens> const& value) override;

  void WriteRecordWithDynamicNDArraysImpl(test_model::RecordWithDynamicNDArrays const& value) override;

  private:
};

// HDF5 reader for the DynamicNDArrays protocol.
class DynamicNDArraysReader : public test_model::DynamicNDArraysReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  DynamicNDArraysReader(std::string path);

  void ReadIntsImpl(yardl::DynamicNDArray<int32_t>& value) override;

  void ReadSimpleRecordArrayImpl(yardl::DynamicNDArray<test_model::SimpleRecord>& value) override;

  void ReadRecordWithVlensArrayImpl(yardl::DynamicNDArray<test_model::RecordWithVlens>& value) override;

  void ReadRecordWithDynamicNDArraysImpl(test_model::RecordWithDynamicNDArrays& value) override;

  private:
};
// HDF5 writer for the Maps protocol.
class MapsWriter : public test_model::MapsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  MapsWriter(std::string path);

  protected:
  void WriteStringToIntImpl(std::unordered_map<std::string, int32_t> const& value) override;

  void WriteStringToUnionImpl(std::unordered_map<std::string, std::variant<std::string, int32_t>> const& value) override;

  void WriteAliasedGenericImpl(test_model::AliasedMap<std::string, int32_t> const& value) override;

  private:
};

// HDF5 reader for the Maps protocol.
class MapsReader : public test_model::MapsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  MapsReader(std::string path);

  void ReadStringToIntImpl(std::unordered_map<std::string, int32_t>& value) override;

  void ReadStringToUnionImpl(std::unordered_map<std::string, std::variant<std::string, int32_t>>& value) override;

  void ReadAliasedGenericImpl(test_model::AliasedMap<std::string, int32_t>& value) override;

  private:
};
// HDF5 writer for the Unions protocol.
class UnionsWriter : public test_model::UnionsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  UnionsWriter(std::string path);

  protected:
  void WriteIntOrSimpleRecordImpl(std::variant<int32_t, test_model::SimpleRecord> const& value) override;

  void WriteIntOrRecordWithVlensImpl(std::variant<int32_t, test_model::RecordWithVlens> const& value) override;

  void WriteMonosotateOrIntOrSimpleRecordImpl(std::variant<std::monostate, int32_t, test_model::SimpleRecord> const& value) override;

  void WriteRecordWithUnionsImpl(test_model::RecordWithUnions const& value) override;

  private:
};

// HDF5 reader for the Unions protocol.
class UnionsReader : public test_model::UnionsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  UnionsReader(std::string path);

  void ReadIntOrSimpleRecordImpl(std::variant<int32_t, test_model::SimpleRecord>& value) override;

  void ReadIntOrRecordWithVlensImpl(std::variant<int32_t, test_model::RecordWithVlens>& value) override;

  void ReadMonosotateOrIntOrSimpleRecordImpl(std::variant<std::monostate, int32_t, test_model::SimpleRecord>& value) override;

  void ReadRecordWithUnionsImpl(test_model::RecordWithUnions& value) override;

  private:
};
// HDF5 writer for the StreamsOfUnions protocol.
class StreamsOfUnionsWriter : public test_model::StreamsOfUnionsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  StreamsOfUnionsWriter(std::string path);

  protected:
  void WriteIntOrSimpleRecordImpl(std::variant<int32_t, test_model::SimpleRecord> const& value) override;

  void EndIntOrSimpleRecordImpl() override;

  void WriteNullableIntOrSimpleRecordImpl(std::variant<std::monostate, int32_t, test_model::SimpleRecord> const& value) override;

  void EndNullableIntOrSimpleRecordImpl() override;

  public:
  void Flush() override;

  private:
  std::unique_ptr<yardl::hdf5::UnionDatasetWriter<2>> intOrSimpleRecord_dataset_state_;
  std::unique_ptr<yardl::hdf5::UnionDatasetWriter<2>> nullableIntOrSimpleRecord_dataset_state_;
};

// HDF5 reader for the StreamsOfUnions protocol.
class StreamsOfUnionsReader : public test_model::StreamsOfUnionsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  StreamsOfUnionsReader(std::string path);

  bool ReadIntOrSimpleRecordImpl(std::variant<int32_t, test_model::SimpleRecord>& value) override;

  bool ReadNullableIntOrSimpleRecordImpl(std::variant<std::monostate, int32_t, test_model::SimpleRecord>& value) override;

  private:
  std::unique_ptr<yardl::hdf5::UnionDatasetReader<2>> intOrSimpleRecord_dataset_state_;
  std::unique_ptr<yardl::hdf5::UnionDatasetReader<2>> nullableIntOrSimpleRecord_dataset_state_;
};
// HDF5 writer for the Enums protocol.
class EnumsWriter : public test_model::EnumsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  EnumsWriter(std::string path);

  protected:
  void WriteSingleImpl(test_model::Fruits const& value) override;

  void WriteVecImpl(std::vector<test_model::Fruits> const& value) override;

  void WriteSizeImpl(test_model::SizeBasedEnum const& value) override;

  private:
};

// HDF5 reader for the Enums protocol.
class EnumsReader : public test_model::EnumsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  EnumsReader(std::string path);

  void ReadSingleImpl(test_model::Fruits& value) override;

  void ReadVecImpl(std::vector<test_model::Fruits>& value) override;

  void ReadSizeImpl(test_model::SizeBasedEnum& value) override;

  private:
};
// HDF5 writer for the Flags protocol.
class FlagsWriter : public test_model::FlagsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  FlagsWriter(std::string path);

  protected:
  void WriteDaysImpl(test_model::DaysOfWeek const& value) override;

  void WriteDaysImpl(std::vector<test_model::DaysOfWeek> const& values) override;

  void EndDaysImpl() override;

  void WriteFormatsImpl(test_model::TextFormat const& value) override;

  void WriteFormatsImpl(std::vector<test_model::TextFormat> const& values) override;

  void EndFormatsImpl() override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetWriter> days_dataset_state_;
  std::unique_ptr<yardl::hdf5::DatasetWriter> formats_dataset_state_;
};

// HDF5 reader for the Flags protocol.
class FlagsReader : public test_model::FlagsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  FlagsReader(std::string path);

  bool ReadDaysImpl(test_model::DaysOfWeek& value) override;

  bool ReadDaysImpl(std::vector<test_model::DaysOfWeek>& values) override;

  bool ReadFormatsImpl(test_model::TextFormat& value) override;

  bool ReadFormatsImpl(std::vector<test_model::TextFormat>& values) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetReader> days_dataset_state_;
  std::unique_ptr<yardl::hdf5::DatasetReader> formats_dataset_state_;
};
// HDF5 writer for the StateTest protocol.
class StateTestWriter : public test_model::StateTestWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  StateTestWriter(std::string path);

  protected:
  void WriteAnIntImpl(int32_t const& value) override;

  void WriteAStreamImpl(int32_t const& value) override;

  void WriteAStreamImpl(std::vector<int32_t> const& values) override;

  void EndAStreamImpl() override;

  void WriteAnotherIntImpl(int32_t const& value) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetWriter> aStream_dataset_state_;
};

// HDF5 reader for the StateTest protocol.
class StateTestReader : public test_model::StateTestReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  StateTestReader(std::string path);

  void ReadAnIntImpl(int32_t& value) override;

  bool ReadAStreamImpl(int32_t& value) override;

  bool ReadAStreamImpl(std::vector<int32_t>& values) override;

  void ReadAnotherIntImpl(int32_t& value) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetReader> aStream_dataset_state_;
};
// HDF5 writer for the SimpleGenerics protocol.
class SimpleGenericsWriter : public test_model::SimpleGenericsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  SimpleGenericsWriter(std::string path);

  protected:
  void WriteFloatImageImpl(test_model::Image<float> const& value) override;

  void WriteIntImageImpl(test_model::Image<int32_t> const& value) override;

  void WriteIntImageAlternateSyntaxImpl(test_model::Image<int32_t> const& value) override;

  void WriteStringImageImpl(test_model::Image<std::string> const& value) override;

  void WriteIntFloatTupleImpl(test_model::MyTuple<int32_t, float> const& value) override;

  void WriteFloatFloatTupleImpl(test_model::MyTuple<float, float> const& value) override;

  void WriteIntFloatTupleAlternateSyntaxImpl(test_model::MyTuple<int32_t, float> const& value) override;

  void WriteIntStringTupleImpl(test_model::MyTuple<int32_t, std::string> const& value) override;

  void WriteStreamOfTypeVariantsImpl(std::variant<test_model::Image<float>, test_model::Image<double>> const& value) override;

  void EndStreamOfTypeVariantsImpl() override;

  public:
  void Flush() override;

  private:
  std::unique_ptr<yardl::hdf5::UnionDatasetWriter<2>> streamOfTypeVariants_dataset_state_;
};

// HDF5 reader for the SimpleGenerics protocol.
class SimpleGenericsReader : public test_model::SimpleGenericsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  SimpleGenericsReader(std::string path);

  void ReadFloatImageImpl(test_model::Image<float>& value) override;

  void ReadIntImageImpl(test_model::Image<int32_t>& value) override;

  void ReadIntImageAlternateSyntaxImpl(test_model::Image<int32_t>& value) override;

  void ReadStringImageImpl(test_model::Image<std::string>& value) override;

  void ReadIntFloatTupleImpl(test_model::MyTuple<int32_t, float>& value) override;

  void ReadFloatFloatTupleImpl(test_model::MyTuple<float, float>& value) override;

  void ReadIntFloatTupleAlternateSyntaxImpl(test_model::MyTuple<int32_t, float>& value) override;

  void ReadIntStringTupleImpl(test_model::MyTuple<int32_t, std::string>& value) override;

  bool ReadStreamOfTypeVariantsImpl(std::variant<test_model::Image<float>, test_model::Image<double>>& value) override;

  private:
  std::unique_ptr<yardl::hdf5::UnionDatasetReader<2>> streamOfTypeVariants_dataset_state_;
};
// HDF5 writer for the AdvancedGenerics protocol.
class AdvancedGenericsWriter : public test_model::AdvancedGenericsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  AdvancedGenericsWriter(std::string path);

  protected:
  void WriteFloatImageImageImpl(test_model::Image<test_model::Image<float>> const& value) override;

  void WriteGenericRecord1Impl(test_model::GenericRecord<int32_t, std::string> const& value) override;

  void WriteTupleOfOptionalsImpl(test_model::MyTuple<std::optional<int32_t>, std::optional<std::string>> const& value) override;

  void WriteTupleOfOptionalsAlternateSyntaxImpl(test_model::MyTuple<std::optional<int32_t>, std::optional<std::string>> const& value) override;

  void WriteTupleOfVectorsImpl(test_model::MyTuple<std::vector<int32_t>, std::vector<float>> const& value) override;

  private:
};

// HDF5 reader for the AdvancedGenerics protocol.
class AdvancedGenericsReader : public test_model::AdvancedGenericsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  AdvancedGenericsReader(std::string path);

  void ReadFloatImageImageImpl(test_model::Image<test_model::Image<float>>& value) override;

  void ReadGenericRecord1Impl(test_model::GenericRecord<int32_t, std::string>& value) override;

  void ReadTupleOfOptionalsImpl(test_model::MyTuple<std::optional<int32_t>, std::optional<std::string>>& value) override;

  void ReadTupleOfOptionalsAlternateSyntaxImpl(test_model::MyTuple<std::optional<int32_t>, std::optional<std::string>>& value) override;

  void ReadTupleOfVectorsImpl(test_model::MyTuple<std::vector<int32_t>, std::vector<float>>& value) override;

  private:
};
// HDF5 writer for the Aliases protocol.
class AliasesWriter : public test_model::AliasesWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  AliasesWriter(std::string path);

  protected:
  void WriteAliasedStringImpl(test_model::AliasedString const& value) override;

  void WriteAliasedEnumImpl(test_model::AliasedEnum const& value) override;

  void WriteAliasedOpenGenericImpl(test_model::AliasedOpenGeneric<test_model::AliasedString, test_model::AliasedEnum> const& value) override;

  void WriteAliasedClosedGenericImpl(test_model::AliasedClosedGeneric const& value) override;

  void WriteAliasedOptionalImpl(test_model::AliasedOptional const& value) override;

  void WriteAliasedGenericOptionalImpl(test_model::AliasedGenericOptional<float> const& value) override;

  void WriteAliasedGenericUnion2Impl(test_model::AliasedGenericUnion2<test_model::AliasedString, test_model::AliasedEnum> const& value) override;

  void WriteAliasedGenericVectorImpl(test_model::AliasedGenericVector<float> const& value) override;

  void WriteAliasedGenericFixedVectorImpl(test_model::AliasedGenericFixedVector<float> const& value) override;

  void WriteStreamOfAliasedGenericUnion2Impl(test_model::AliasedGenericUnion2<test_model::AliasedString, test_model::AliasedEnum> const& value) override;

  void EndStreamOfAliasedGenericUnion2Impl() override;

  public:
  void Flush() override;

  private:
  std::unique_ptr<yardl::hdf5::UnionDatasetWriter<2>> streamOfAliasedGenericUnion2_dataset_state_;
};

// HDF5 reader for the Aliases protocol.
class AliasesReader : public test_model::AliasesReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  AliasesReader(std::string path);

  void ReadAliasedStringImpl(test_model::AliasedString& value) override;

  void ReadAliasedEnumImpl(test_model::AliasedEnum& value) override;

  void ReadAliasedOpenGenericImpl(test_model::AliasedOpenGeneric<test_model::AliasedString, test_model::AliasedEnum>& value) override;

  void ReadAliasedClosedGenericImpl(test_model::AliasedClosedGeneric& value) override;

  void ReadAliasedOptionalImpl(test_model::AliasedOptional& value) override;

  void ReadAliasedGenericOptionalImpl(test_model::AliasedGenericOptional<float>& value) override;

  void ReadAliasedGenericUnion2Impl(test_model::AliasedGenericUnion2<test_model::AliasedString, test_model::AliasedEnum>& value) override;

  void ReadAliasedGenericVectorImpl(test_model::AliasedGenericVector<float>& value) override;

  void ReadAliasedGenericFixedVectorImpl(test_model::AliasedGenericFixedVector<float>& value) override;

  bool ReadStreamOfAliasedGenericUnion2Impl(test_model::AliasedGenericUnion2<test_model::AliasedString, test_model::AliasedEnum>& value) override;

  private:
  std::unique_ptr<yardl::hdf5::UnionDatasetReader<2>> streamOfAliasedGenericUnion2_dataset_state_;
};
// HDF5 writer for the StreamsOfAliasedUnions protocol.
class StreamsOfAliasedUnionsWriter : public test_model::StreamsOfAliasedUnionsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  StreamsOfAliasedUnionsWriter(std::string path);

  protected:
  void WriteIntOrSimpleRecordImpl(test_model::AliasedIntOrSimpleRecord const& value) override;

  void EndIntOrSimpleRecordImpl() override;

  void WriteNullableIntOrSimpleRecordImpl(test_model::AliasedNullableIntSimpleRecord const& value) override;

  void EndNullableIntOrSimpleRecordImpl() override;

  public:
  void Flush() override;

  private:
  std::unique_ptr<yardl::hdf5::UnionDatasetWriter<2>> intOrSimpleRecord_dataset_state_;
  std::unique_ptr<yardl::hdf5::UnionDatasetWriter<2>> nullableIntOrSimpleRecord_dataset_state_;
};

// HDF5 reader for the StreamsOfAliasedUnions protocol.
class StreamsOfAliasedUnionsReader : public test_model::StreamsOfAliasedUnionsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  StreamsOfAliasedUnionsReader(std::string path);

  bool ReadIntOrSimpleRecordImpl(test_model::AliasedIntOrSimpleRecord& value) override;

  bool ReadNullableIntOrSimpleRecordImpl(test_model::AliasedNullableIntSimpleRecord& value) override;

  private:
  std::unique_ptr<yardl::hdf5::UnionDatasetReader<2>> intOrSimpleRecord_dataset_state_;
  std::unique_ptr<yardl::hdf5::UnionDatasetReader<2>> nullableIntOrSimpleRecord_dataset_state_;
};
// HDF5 writer for the ProtocolWithComputedFields protocol.
class ProtocolWithComputedFieldsWriter : public test_model::ProtocolWithComputedFieldsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  ProtocolWithComputedFieldsWriter(std::string path);

  protected:
  void WriteRecordWithComputedFieldsImpl(test_model::RecordWithComputedFields const& value) override;

  private:
};

// HDF5 reader for the ProtocolWithComputedFields protocol.
class ProtocolWithComputedFieldsReader : public test_model::ProtocolWithComputedFieldsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  ProtocolWithComputedFieldsReader(std::string path);

  void ReadRecordWithComputedFieldsImpl(test_model::RecordWithComputedFields& value) override;

  private:
};
// HDF5 writer for the ProtocolWithKeywordSteps protocol.
class ProtocolWithKeywordStepsWriter : public test_model::ProtocolWithKeywordStepsWriterBase, public yardl::hdf5::Hdf5Writer {
  public:
  ProtocolWithKeywordStepsWriter(std::string path);

  protected:
  void WriteIntImpl(test_model::RecordWithKeywordFields const& value) override;

  void WriteIntImpl(std::vector<test_model::RecordWithKeywordFields> const& values) override;

  void EndIntImpl() override;

  void WriteFloatImpl(test_model::EnumWithKeywordSymbols const& value) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetWriter> int_dataset_state_;
};

// HDF5 reader for the ProtocolWithKeywordSteps protocol.
class ProtocolWithKeywordStepsReader : public test_model::ProtocolWithKeywordStepsReaderBase, public yardl::hdf5::Hdf5Reader {
  public:
  ProtocolWithKeywordStepsReader(std::string path);

  bool ReadIntImpl(test_model::RecordWithKeywordFields& value) override;

  bool ReadIntImpl(std::vector<test_model::RecordWithKeywordFields>& values) override;

  void ReadFloatImpl(test_model::EnumWithKeywordSymbols& value) override;

  private:
  std::unique_ptr<yardl::hdf5::DatasetReader> int_dataset_state_;
};
}
