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
    AdvancedGenericsReaderBase,
    AdvancedGenericsWriterBase,
    AliasesReaderBase,
    AliasesWriterBase,
    BenchmarkFloat256x256ReaderBase,
    BenchmarkFloat256x256WriterBase,
    BenchmarkFloatVlenReaderBase,
    BenchmarkFloatVlenWriterBase,
    BenchmarkInt256x256ReaderBase,
    BenchmarkInt256x256WriterBase,
    BenchmarkSimpleMrdReaderBase,
    BenchmarkSimpleMrdWriterBase,
    BenchmarkSmallRecordReaderBase,
    BenchmarkSmallRecordWithOptionalsReaderBase,
    BenchmarkSmallRecordWithOptionalsWriterBase,
    BenchmarkSmallRecordWriterBase,
    ComplexArraysReaderBase,
    ComplexArraysWriterBase,
    DynamicNDArraysReaderBase,
    DynamicNDArraysWriterBase,
    EnumsReaderBase,
    EnumsWriterBase,
    FixedArraysReaderBase,
    FixedArraysWriterBase,
    FixedVectorsReaderBase,
    FixedVectorsWriterBase,
    FlagsReaderBase,
    FlagsWriterBase,
    MapsReaderBase,
    MapsWriterBase,
    MultiDArraysReaderBase,
    MultiDArraysWriterBase,
    NDArraysReaderBase,
    NDArraysSingleDimensionReaderBase,
    NDArraysSingleDimensionWriterBase,
    NDArraysWriterBase,
    NestedRecordsReaderBase,
    NestedRecordsWriterBase,
    OptionalVectorsReaderBase,
    OptionalVectorsWriterBase,
    ProtocolWithComputedFieldsReaderBase,
    ProtocolWithComputedFieldsWriterBase,
    ProtocolWithKeywordStepsReaderBase,
    ProtocolWithKeywordStepsWriterBase,
    ProtocolWithOptionalDateReaderBase,
    ProtocolWithOptionalDateWriterBase,
    ScalarOptionalsReaderBase,
    ScalarOptionalsWriterBase,
    ScalarsReaderBase,
    ScalarsWriterBase,
    SimpleGenericsReaderBase,
    SimpleGenericsWriterBase,
    StateTestReaderBase,
    StateTestWriterBase,
    StreamsOfAliasedUnionsReaderBase,
    StreamsOfAliasedUnionsWriterBase,
    StreamsOfUnionsReaderBase,
    StreamsOfUnionsWriterBase,
    StreamsReaderBase,
    StreamsWriterBase,
    StringsReaderBase,
    StringsWriterBase,
    SubarraysInRecordsReaderBase,
    SubarraysInRecordsWriterBase,
    SubarraysReaderBase,
    SubarraysWriterBase,
    UnionsReaderBase,
    UnionsWriterBase,
    VlensReaderBase,
    VlensWriterBase,
)
from .binary import (
    BinaryAdvancedGenericsReader,
    BinaryAdvancedGenericsWriter,
    BinaryAliasesReader,
    BinaryAliasesWriter,
    BinaryBenchmarkFloat256x256Reader,
    BinaryBenchmarkFloat256x256Writer,
    BinaryBenchmarkFloatVlenReader,
    BinaryBenchmarkFloatVlenWriter,
    BinaryBenchmarkInt256x256Reader,
    BinaryBenchmarkInt256x256Writer,
    BinaryBenchmarkSimpleMrdReader,
    BinaryBenchmarkSimpleMrdWriter,
    BinaryBenchmarkSmallRecordReader,
    BinaryBenchmarkSmallRecordWithOptionalsReader,
    BinaryBenchmarkSmallRecordWithOptionalsWriter,
    BinaryBenchmarkSmallRecordWriter,
    BinaryComplexArraysReader,
    BinaryComplexArraysWriter,
    BinaryDynamicNDArraysReader,
    BinaryDynamicNDArraysWriter,
    BinaryEnumsReader,
    BinaryEnumsWriter,
    BinaryFixedArraysReader,
    BinaryFixedArraysWriter,
    BinaryFixedVectorsReader,
    BinaryFixedVectorsWriter,
    BinaryFlagsReader,
    BinaryFlagsWriter,
    BinaryMapsReader,
    BinaryMapsWriter,
    BinaryMultiDArraysReader,
    BinaryMultiDArraysWriter,
    BinaryNDArraysReader,
    BinaryNDArraysSingleDimensionReader,
    BinaryNDArraysSingleDimensionWriter,
    BinaryNDArraysWriter,
    BinaryNestedRecordsReader,
    BinaryNestedRecordsWriter,
    BinaryOptionalVectorsReader,
    BinaryOptionalVectorsWriter,
    BinaryProtocolWithComputedFieldsReader,
    BinaryProtocolWithComputedFieldsWriter,
    BinaryProtocolWithKeywordStepsReader,
    BinaryProtocolWithKeywordStepsWriter,
    BinaryProtocolWithOptionalDateReader,
    BinaryProtocolWithOptionalDateWriter,
    BinaryScalarOptionalsReader,
    BinaryScalarOptionalsWriter,
    BinaryScalarsReader,
    BinaryScalarsWriter,
    BinarySimpleGenericsReader,
    BinarySimpleGenericsWriter,
    BinaryStateTestReader,
    BinaryStateTestWriter,
    BinaryStreamsOfAliasedUnionsReader,
    BinaryStreamsOfAliasedUnionsWriter,
    BinaryStreamsOfUnionsReader,
    BinaryStreamsOfUnionsWriter,
    BinaryStreamsReader,
    BinaryStreamsWriter,
    BinaryStringsReader,
    BinaryStringsWriter,
    BinarySubarraysInRecordsReader,
    BinarySubarraysInRecordsWriter,
    BinarySubarraysReader,
    BinarySubarraysWriter,
    BinaryUnionsReader,
    BinaryUnionsWriter,
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
