import numpy as np
import pytest
import test_model as tm


def test_field_access():
    r = tm.RecordWithComputedFields()

    r.int_field = 42
    assert r.access_int_field() == 42

    r.string_field = "hello"
    assert r.access_string_field() == "hello"

    r.tuple_field = tm.MyTuple(v1=1, v2=2)
    assert r.access_tuple_field() == r.tuple_field
    assert r.access_nested_tuple_field() == r.tuple_field.v2

    r.array_field = np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32)
    assert r.access_array_field() is r.array_field
    assert r.access_array_field_element() == r.array_field[0, 1]
    assert r.access_array_field_element_by_name() == r.array_field[0, 1]

    assert r.access_other_computed_field() == r.access_int_field()

    r.vector_of_vectors_field = [[1, 2, 3], [4, 5, 6]]
    assert r.access_vector_of_vectors_field() == r.vector_of_vectors_field[1][2]

    r.map_field = {"hello": "world", "world": "bye"}
    assert r.access_map() is r.map_field
    assert r.access_map_entry() == "world"
    assert r.access_map_entry_with_computed_field() == "world"
    assert r.access_map_entry_with_computed_field_nested() == "bye"
    with pytest.raises(KeyError):
        r.access_missing_map_entry()


def test_literals():
    r = tm.RecordWithComputedFields()
    assert r.int_literal() == 42
    assert r.large_negative_int64_literal() == -0x4000000000000000
    assert r.large_u_int64_literal() == 0x8000000000000000
    assert r.string_literal() == "hello"


def test_dimension_index():
    r = tm.RecordWithComputedFields()
    assert r.array_dimension_x_index() == 0
    assert r.array_dimension_y_index() == 1

    r.string_field = "y"
    assert r.array_dimension_index_from_string_field() == 1
    with pytest.raises(KeyError):
        r.string_field = "missing"
        r.array_dimension_index_from_string_field()


def test_dimension_count():
    r = tm.RecordWithComputedFields()
    assert r.array_dimension_count() == 2

    r.dynamic_array_field = np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32)
    assert r.dynamic_array_dimension_count() == 2
    r.dynamic_array_field = np.array([1, 2, 3], dtype=np.int32)
    assert r.dynamic_array_dimension_count() == 1


def test_vector_size():
    r = tm.RecordWithComputedFields()
    assert r.vector_size() == 0
    r.vector_field = [1, 2, 3, 4]
    assert r.vector_size() == 4

    assert r.fixed_vector_size() == 3


def test_map_size():
    r = tm.RecordWithComputedFields()
    assert r.map_size() == 0
    r.map_field = {"hello": "bonjour", "world": "monde"}
    assert r.map_size() == 2


def test_array_size():
    r = tm.RecordWithComputedFields()
    r.array_field = np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32)

    assert r.array_size() == 6
    assert r.array_x_size() == 2
    assert r.array_y_size() == 3
    assert r.array_0_size() == 2
    assert r.array_1_size() == 3

    assert r.array_size_from_int_field() == 2
    r.int_field = 1
    assert r.array_size_from_int_field() == 3

    r.string_field = "x"
    assert r.array_size_from_string_field() == 2
    r.string_field = "y"
    assert r.array_size_from_string_field() == 3

    with pytest.raises(KeyError):
        r.string_field = "missing"
        r.array_size_from_string_field()

    r.tuple_field.v1 = 1
    assert r.array_size_from_nested_int_field() == 3
    assert r.fixed_array_size() == r.fixed_array_field.size
    assert r.fixed_array_x_size() == r.fixed_array_field.shape[0]
    assert r.fixed_array_0_size() == r.fixed_array_field.shape[0]

    r.array_field_map_dimensions = np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32)
    assert r.array_field_map_dimensions_x_size() == 2


def test_switch_expression():
    r = tm.RecordWithComputedFields()
    r.optional_named_array = np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32)
    assert r.optional_named_array_length() == 6
    assert r.optional_named_array_length_with_discard() == 6

    r.optional_named_array = None
    assert r.optional_named_array_length() == 0
    assert r.optional_named_array_length_with_discard() == 0

    r.int_float_union = tm.Int32OrFloat32.Int32(42)
    assert r.int_float_union_as_float() == 42.0
    r.int_float_union = tm.Int32OrFloat32.Float32(42.9)
    assert r.int_float_union_as_float() == 42.9

    r.nullable_int_float_union = None
    assert r.nullable_int_float_union_string() == "null"
    r.nullable_int_float_union = tm.Int32OrFloat32.Int32(42)
    assert r.nullable_int_float_union_string() == "int"
    r.nullable_int_float_union = tm.Int32OrFloat32.Float32(42.9)
    assert r.nullable_int_float_union_string() == "float"

    r.union_with_nested_generic_union = tm.IntOrGenericRecordWithComputedFields.Int(42)
    assert r.nested_switch() == -1
    assert r.use_nested_computed_field() == -1

    r.union_with_nested_generic_union = (
        tm.IntOrGenericRecordWithComputedFields.GenericRecordWithComputedFields(
            tm.basic_types.GenericRecordWithComputedFields(
                f1=tm.basic_types.T0OrT1[str, tm.Float32].T0("hi")
            )
        )
    )

    assert r.nested_switch() == 10
    assert r.use_nested_computed_field() == 0
    r.union_with_nested_generic_union = (
        tm.IntOrGenericRecordWithComputedFields.GenericRecordWithComputedFields(
            tm.basic_types.GenericRecordWithComputedFields(
                f1=tm.basic_types.T0OrT1[str, tm.Float32].T1(42.9)
            )
        )
    )
    assert r.nested_switch() == 20
    assert r.use_nested_computed_field() == 1

    r.int_field = 10
    assert r.switch_over_single_value() == 10

    gr = tm.basic_types.GenericRecordWithComputedFields(
        f1=tm.basic_types.T0OrT1[tm.Int32, tm.Float32].T0(42)
    )
    assert gr.type_index() == 0
    gr.f1 = tm.basic_types.T0OrT1[tm.Int32, tm.Float32].T1(42.9)
    assert gr.type_index() == 1


def test_arithmetic():
    r = tm.RecordWithComputedFields()
    assert r.arithmetic_1() == 3
    assert r.arithmetic_2() == 11
    assert r.arithmetic_3() == 13

    r.array_field = np.array([[1, 2, 3], [4, 5, 6]], dtype=np.int32)
    r.int_field = 1
    assert r.arithmetic_4() == 5
    assert r.arithmetic_5() == 3

    assert r.arithmetic_6() == 3

    assert r.arithmetic_7() == 49.0
    assert isinstance(r.arithmetic_7(), float)

    r.complexfloat32_field = 2.0 + 3.0j
    assert r.arithmetic8() == 6.0 + 9.0j

    assert r.arithmetic_9() == 2.2

    assert r.arithmetic_10() == 1e10 + 9e9


def test_casting():
    r = tm.RecordWithComputedFields()
    r.int_field = 42
    assert r.cast_int_to_float() == 42.0
    assert isinstance(r.cast_int_to_float(), float)

    r.float32_field = 42.9
    assert r.cast_float_to_int() == 42
    assert isinstance(r.cast_float_to_int(), int)

    assert r.cast_power() == 49
    assert isinstance(r.cast_power(), int)

    r.complexfloat32_field = 2.0 + 3.0j
    r.complexfloat64_field = 2.0 + 3.0j
    assert r.cast_complex32_to_complex64() == 2.0 + 3.0j
    assert r.cast_complex64_to_complex32() == 2.0 + 3.0j

    assert r.cast_float_to_complex() == 66.6 + 0.0j
    assert isinstance(r.cast_float_to_complex(), complex)
