add_executable(test_only_binary only_binary.cc)

target_link_libraries(
  test_only_binary
  only_binary_generated
)

if(EXISTS ${CMAKE_CURRENT_SOURCE_DIR}/generated)
  add_subdirectory(generated)
endif()
