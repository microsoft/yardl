# This file was generated by the "yardl" tool. DO NOT EDIT.

# pyright: reportUnusedImport=false
from typing import Tuple as _Tuple
import re as _re
import numpy as _np

_MIN_NUMPY_VERSION = (1, 22, 0)

def _parse_version(version: str) -> _Tuple[int, ...]:
    try:
        return tuple(map(int, version.split(".")))
    except ValueError:
        # ignore any prerelease suffix
        version = _re.sub(r"[^0-9.]", "", version)
        return tuple(map(int, version.split(".")))

if _parse_version(_np.__version__) < _MIN_NUMPY_VERSION:
    raise ImportError(f"Your installed numpy version is {_np.__version__}, but version >= {'.'.join(str(i) for i in _MIN_NUMPY_VERSION)} is required.")

from .yardl_types import *
from . import tuples
from . import basic_types
from . import image
from .types import (
    AcquisitionOrImage,
    AliasedClosedGeneric,
    AliasedEnum,
    AliasedGenericDynamicArray,
    AliasedGenericFixedArray,
    AliasedGenericFixedVector,
    AliasedGenericOptional,
    AliasedGenericRank2Array,
    AliasedGenericUnion2,
    AliasedGenericVector,
    AliasedIntOrAliasedSimpleRecord,
    AliasedIntOrSimpleRecord,
    AliasedMap,
    AliasedMultiGenericOptional,
    AliasedNullableIntSimpleRecord,
    AliasedOpenGeneric,
    AliasedOptional,
    AliasedSimpleRecord,
    AliasedString,
    AliasedTuple,
    AliasedVectorOfGenericRecords,
    ArrayOrScalar,
    ArrayWithKeywordDimensionNames,
    DaysOfWeek,
    EnumWithKeywordSymbols,
    Fruits,
    GenericRecord,
    GenericUnion3,
    GenericUnion3Alternate,
    GenericUnionWithRepeatedTypeParameters,
    Image,
    ImageFloatOrImageDouble,
    Int32OrFloat32,
    Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray,
    Int32OrRecordWithVlens,
    Int32OrSimpleRecord,
    Int64Enum,
    IntArray,
    IntFixedArray,
    IntOrGenericRecordWithComputedFields,
    IntRank2Array,
    MapOrScalar,
    MyTuple,
    NamedFixedNDArray,
    NamedNDArray,
    RecordContainingGenericRecords,
    RecordContainingNestedGenericRecords,
    RecordNotUsedInProtocol,
    RecordWithAliasedGenerics,
    RecordWithAliasedOptionalGenericField,
    RecordWithAliasedOptionalGenericUnionField,
    RecordWithArrays,
    RecordWithArraysSimpleSyntax,
    RecordWithComputedFields,
    RecordWithDynamicNDArrays,
    RecordWithEnums,
    RecordWithFixedArrays,
    RecordWithFixedCollections,
    RecordWithFixedVectors,
    RecordWithFloatArrays,
    RecordWithGenericArrays,
    RecordWithGenericFixedVectors,
    RecordWithGenericMaps,
    RecordWithGenericVectorOfRecords,
    RecordWithGenericVectors,
    RecordWithIntVectors,
    RecordWithKeywordFields,
    RecordWithMaps,
    RecordWithNDArrays,
    RecordWithNDArraysSingleDimension,
    RecordWithNamedFixedArrays,
    RecordWithNoDefaultEnum,
    RecordWithOptionalDate,
    RecordWithOptionalFields,
    RecordWithOptionalGenericField,
    RecordWithOptionalGenericUnionField,
    RecordWithOptionalVector,
    RecordWithPrimitiveAliases,
    RecordWithPrimitives,
    RecordWithStrings,
    RecordWithUnionsOfContainers,
    RecordWithVectorOfTimes,
    RecordWithVectors,
    RecordWithVlenCollections,
    RecordWithVlens,
    RecordWithVlensFixedArray,
    SimpleAcquisition,
    SimpleEncodingCounters,
    SimpleRecord,
    SimpleRecordFixedArray,
    SizeBasedEnum,
    SmallBenchmarkRecord,
    StringOrInt32,
    TextFormat,
    TupleWithRecords,
    UInt64Enum,
    UOrV,
    UnionOfContainerRecords,
    VectorOfGenericRecords,
    VectorOrScalar,
    get_dtype,
)
from .protocols import (
    AdvancedGenericsIndexedReaderBase,
    AdvancedGenericsReaderBase,
    AdvancedGenericsWriterBase,
    AliasesIndexedReaderBase,
    AliasesReaderBase,
    AliasesWriterBase,
    BenchmarkFloat256x256IndexedReaderBase,
    BenchmarkFloat256x256ReaderBase,
    BenchmarkFloat256x256WriterBase,
    BenchmarkFloatVlenIndexedReaderBase,
    BenchmarkFloatVlenReaderBase,
    BenchmarkFloatVlenWriterBase,
    BenchmarkInt256x256IndexedReaderBase,
    BenchmarkInt256x256ReaderBase,
    BenchmarkInt256x256WriterBase,
    BenchmarkSimpleMrdIndexedReaderBase,
    BenchmarkSimpleMrdReaderBase,
    BenchmarkSimpleMrdWriterBase,
    BenchmarkSmallRecordIndexedReaderBase,
    BenchmarkSmallRecordReaderBase,
    BenchmarkSmallRecordWithOptionalsIndexedReaderBase,
    BenchmarkSmallRecordWithOptionalsReaderBase,
    BenchmarkSmallRecordWithOptionalsWriterBase,
    BenchmarkSmallRecordWriterBase,
    ComplexArraysIndexedReaderBase,
    ComplexArraysReaderBase,
    ComplexArraysWriterBase,
    DynamicNDArraysIndexedReaderBase,
    DynamicNDArraysReaderBase,
    DynamicNDArraysWriterBase,
    EnumsIndexedReaderBase,
    EnumsReaderBase,
    EnumsWriterBase,
    FixedArraysIndexedReaderBase,
    FixedArraysReaderBase,
    FixedArraysWriterBase,
    FixedVectorsIndexedReaderBase,
    FixedVectorsReaderBase,
    FixedVectorsWriterBase,
    FlagsIndexedReaderBase,
    FlagsReaderBase,
    FlagsWriterBase,
    MapsIndexedReaderBase,
    MapsReaderBase,
    MapsWriterBase,
    MultiDArraysIndexedReaderBase,
    MultiDArraysReaderBase,
    MultiDArraysWriterBase,
    NDArraysIndexedReaderBase,
    NDArraysReaderBase,
    NDArraysSingleDimensionIndexedReaderBase,
    NDArraysSingleDimensionReaderBase,
    NDArraysSingleDimensionWriterBase,
    NDArraysWriterBase,
    NestedRecordsIndexedReaderBase,
    NestedRecordsReaderBase,
    NestedRecordsWriterBase,
    OptionalVectorsIndexedReaderBase,
    OptionalVectorsReaderBase,
    OptionalVectorsWriterBase,
    ProtocolWithComputedFieldsIndexedReaderBase,
    ProtocolWithComputedFieldsReaderBase,
    ProtocolWithComputedFieldsWriterBase,
    ProtocolWithKeywordStepsIndexedReaderBase,
    ProtocolWithKeywordStepsReaderBase,
    ProtocolWithKeywordStepsWriterBase,
    ProtocolWithOptionalDateIndexedReaderBase,
    ProtocolWithOptionalDateReaderBase,
    ProtocolWithOptionalDateWriterBase,
    ScalarOptionalsIndexedReaderBase,
    ScalarOptionalsReaderBase,
    ScalarOptionalsWriterBase,
    ScalarsIndexedReaderBase,
    ScalarsReaderBase,
    ScalarsWriterBase,
    SimpleGenericsIndexedReaderBase,
    SimpleGenericsReaderBase,
    SimpleGenericsWriterBase,
    StateTestIndexedReaderBase,
    StateTestReaderBase,
    StateTestWriterBase,
    StreamsIndexedReaderBase,
    StreamsOfAliasedUnionsIndexedReaderBase,
    StreamsOfAliasedUnionsReaderBase,
    StreamsOfAliasedUnionsWriterBase,
    StreamsOfUnionsIndexedReaderBase,
    StreamsOfUnionsReaderBase,
    StreamsOfUnionsWriterBase,
    StreamsReaderBase,
    StreamsWriterBase,
    StringsIndexedReaderBase,
    StringsReaderBase,
    StringsWriterBase,
    SubarraysInRecordsIndexedReaderBase,
    SubarraysInRecordsReaderBase,
    SubarraysInRecordsWriterBase,
    SubarraysIndexedReaderBase,
    SubarraysReaderBase,
    SubarraysWriterBase,
    UnionsIndexedReaderBase,
    UnionsReaderBase,
    UnionsWriterBase,
    VlensIndexedReaderBase,
    VlensReaderBase,
    VlensWriterBase,
)
from .binary import (
    BinaryAdvancedGenericsIndexedReader,
    BinaryAdvancedGenericsIndexedWriter,
    BinaryAdvancedGenericsReader,
    BinaryAdvancedGenericsWriter,
    BinaryAliasesIndexedReader,
    BinaryAliasesIndexedWriter,
    BinaryAliasesReader,
    BinaryAliasesWriter,
    BinaryBenchmarkFloat256x256IndexedReader,
    BinaryBenchmarkFloat256x256IndexedWriter,
    BinaryBenchmarkFloat256x256Reader,
    BinaryBenchmarkFloat256x256Writer,
    BinaryBenchmarkFloatVlenIndexedReader,
    BinaryBenchmarkFloatVlenIndexedWriter,
    BinaryBenchmarkFloatVlenReader,
    BinaryBenchmarkFloatVlenWriter,
    BinaryBenchmarkInt256x256IndexedReader,
    BinaryBenchmarkInt256x256IndexedWriter,
    BinaryBenchmarkInt256x256Reader,
    BinaryBenchmarkInt256x256Writer,
    BinaryBenchmarkSimpleMrdIndexedReader,
    BinaryBenchmarkSimpleMrdIndexedWriter,
    BinaryBenchmarkSimpleMrdReader,
    BinaryBenchmarkSimpleMrdWriter,
    BinaryBenchmarkSmallRecordIndexedReader,
    BinaryBenchmarkSmallRecordIndexedWriter,
    BinaryBenchmarkSmallRecordReader,
    BinaryBenchmarkSmallRecordWithOptionalsIndexedReader,
    BinaryBenchmarkSmallRecordWithOptionalsIndexedWriter,
    BinaryBenchmarkSmallRecordWithOptionalsReader,
    BinaryBenchmarkSmallRecordWithOptionalsWriter,
    BinaryBenchmarkSmallRecordWriter,
    BinaryComplexArraysIndexedReader,
    BinaryComplexArraysIndexedWriter,
    BinaryComplexArraysReader,
    BinaryComplexArraysWriter,
    BinaryDynamicNDArraysIndexedReader,
    BinaryDynamicNDArraysIndexedWriter,
    BinaryDynamicNDArraysReader,
    BinaryDynamicNDArraysWriter,
    BinaryEnumsIndexedReader,
    BinaryEnumsIndexedWriter,
    BinaryEnumsReader,
    BinaryEnumsWriter,
    BinaryFixedArraysIndexedReader,
    BinaryFixedArraysIndexedWriter,
    BinaryFixedArraysReader,
    BinaryFixedArraysWriter,
    BinaryFixedVectorsIndexedReader,
    BinaryFixedVectorsIndexedWriter,
    BinaryFixedVectorsReader,
    BinaryFixedVectorsWriter,
    BinaryFlagsIndexedReader,
    BinaryFlagsIndexedWriter,
    BinaryFlagsReader,
    BinaryFlagsWriter,
    BinaryMapsIndexedReader,
    BinaryMapsIndexedWriter,
    BinaryMapsReader,
    BinaryMapsWriter,
    BinaryMultiDArraysIndexedReader,
    BinaryMultiDArraysIndexedWriter,
    BinaryMultiDArraysReader,
    BinaryMultiDArraysWriter,
    BinaryNDArraysIndexedReader,
    BinaryNDArraysIndexedWriter,
    BinaryNDArraysReader,
    BinaryNDArraysSingleDimensionIndexedReader,
    BinaryNDArraysSingleDimensionIndexedWriter,
    BinaryNDArraysSingleDimensionReader,
    BinaryNDArraysSingleDimensionWriter,
    BinaryNDArraysWriter,
    BinaryNestedRecordsIndexedReader,
    BinaryNestedRecordsIndexedWriter,
    BinaryNestedRecordsReader,
    BinaryNestedRecordsWriter,
    BinaryOptionalVectorsIndexedReader,
    BinaryOptionalVectorsIndexedWriter,
    BinaryOptionalVectorsReader,
    BinaryOptionalVectorsWriter,
    BinaryProtocolWithComputedFieldsIndexedReader,
    BinaryProtocolWithComputedFieldsIndexedWriter,
    BinaryProtocolWithComputedFieldsReader,
    BinaryProtocolWithComputedFieldsWriter,
    BinaryProtocolWithKeywordStepsIndexedReader,
    BinaryProtocolWithKeywordStepsIndexedWriter,
    BinaryProtocolWithKeywordStepsReader,
    BinaryProtocolWithKeywordStepsWriter,
    BinaryProtocolWithOptionalDateIndexedReader,
    BinaryProtocolWithOptionalDateIndexedWriter,
    BinaryProtocolWithOptionalDateReader,
    BinaryProtocolWithOptionalDateWriter,
    BinaryScalarOptionalsIndexedReader,
    BinaryScalarOptionalsIndexedWriter,
    BinaryScalarOptionalsReader,
    BinaryScalarOptionalsWriter,
    BinaryScalarsIndexedReader,
    BinaryScalarsIndexedWriter,
    BinaryScalarsReader,
    BinaryScalarsWriter,
    BinarySimpleGenericsIndexedReader,
    BinarySimpleGenericsIndexedWriter,
    BinarySimpleGenericsReader,
    BinarySimpleGenericsWriter,
    BinaryStateTestIndexedReader,
    BinaryStateTestIndexedWriter,
    BinaryStateTestReader,
    BinaryStateTestWriter,
    BinaryStreamsIndexedReader,
    BinaryStreamsIndexedWriter,
    BinaryStreamsOfAliasedUnionsIndexedReader,
    BinaryStreamsOfAliasedUnionsIndexedWriter,
    BinaryStreamsOfAliasedUnionsReader,
    BinaryStreamsOfAliasedUnionsWriter,
    BinaryStreamsOfUnionsIndexedReader,
    BinaryStreamsOfUnionsIndexedWriter,
    BinaryStreamsOfUnionsReader,
    BinaryStreamsOfUnionsWriter,
    BinaryStreamsReader,
    BinaryStreamsWriter,
    BinaryStringsIndexedReader,
    BinaryStringsIndexedWriter,
    BinaryStringsReader,
    BinaryStringsWriter,
    BinarySubarraysInRecordsIndexedReader,
    BinarySubarraysInRecordsIndexedWriter,
    BinarySubarraysInRecordsReader,
    BinarySubarraysInRecordsWriter,
    BinarySubarraysIndexedReader,
    BinarySubarraysIndexedWriter,
    BinarySubarraysReader,
    BinarySubarraysWriter,
    BinaryUnionsIndexedReader,
    BinaryUnionsIndexedWriter,
    BinaryUnionsReader,
    BinaryUnionsWriter,
    BinaryVlensIndexedReader,
    BinaryVlensIndexedWriter,
    BinaryVlensReader,
    BinaryVlensWriter,
)
from .ndjson import (
    NDJsonAdvancedGenericsReader,
    NDJsonAdvancedGenericsWriter,
    NDJsonAliasesReader,
    NDJsonAliasesWriter,
    NDJsonBenchmarkFloat256x256Reader,
    NDJsonBenchmarkFloat256x256Writer,
    NDJsonBenchmarkFloatVlenReader,
    NDJsonBenchmarkFloatVlenWriter,
    NDJsonBenchmarkInt256x256Reader,
    NDJsonBenchmarkInt256x256Writer,
    NDJsonBenchmarkSimpleMrdReader,
    NDJsonBenchmarkSimpleMrdWriter,
    NDJsonBenchmarkSmallRecordReader,
    NDJsonBenchmarkSmallRecordWithOptionalsReader,
    NDJsonBenchmarkSmallRecordWithOptionalsWriter,
    NDJsonBenchmarkSmallRecordWriter,
    NDJsonComplexArraysReader,
    NDJsonComplexArraysWriter,
    NDJsonDynamicNDArraysReader,
    NDJsonDynamicNDArraysWriter,
    NDJsonEnumsReader,
    NDJsonEnumsWriter,
    NDJsonFixedArraysReader,
    NDJsonFixedArraysWriter,
    NDJsonFixedVectorsReader,
    NDJsonFixedVectorsWriter,
    NDJsonFlagsReader,
    NDJsonFlagsWriter,
    NDJsonMapsReader,
    NDJsonMapsWriter,
    NDJsonMultiDArraysReader,
    NDJsonMultiDArraysWriter,
    NDJsonNDArraysReader,
    NDJsonNDArraysSingleDimensionReader,
    NDJsonNDArraysSingleDimensionWriter,
    NDJsonNDArraysWriter,
    NDJsonNestedRecordsReader,
    NDJsonNestedRecordsWriter,
    NDJsonOptionalVectorsReader,
    NDJsonOptionalVectorsWriter,
    NDJsonProtocolWithComputedFieldsReader,
    NDJsonProtocolWithComputedFieldsWriter,
    NDJsonProtocolWithKeywordStepsReader,
    NDJsonProtocolWithKeywordStepsWriter,
    NDJsonProtocolWithOptionalDateReader,
    NDJsonProtocolWithOptionalDateWriter,
    NDJsonScalarOptionalsReader,
    NDJsonScalarOptionalsWriter,
    NDJsonScalarsReader,
    NDJsonScalarsWriter,
    NDJsonSimpleGenericsReader,
    NDJsonSimpleGenericsWriter,
    NDJsonStateTestReader,
    NDJsonStateTestWriter,
    NDJsonStreamsOfAliasedUnionsReader,
    NDJsonStreamsOfAliasedUnionsWriter,
    NDJsonStreamsOfUnionsReader,
    NDJsonStreamsOfUnionsWriter,
    NDJsonStreamsReader,
    NDJsonStreamsWriter,
    NDJsonStringsReader,
    NDJsonStringsWriter,
    NDJsonSubarraysInRecordsReader,
    NDJsonSubarraysInRecordsWriter,
    NDJsonSubarraysReader,
    NDJsonSubarraysWriter,
    NDJsonUnionsReader,
    NDJsonUnionsWriter,
    NDJsonVlensReader,
    NDJsonVlensWriter,
)
