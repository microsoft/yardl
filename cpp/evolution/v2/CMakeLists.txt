add_subdirectory(lib)

add_executable(v2_write
  write.cc
)

target_link_libraries(v2_write
  evo_test_generated
)

add_executable(v2_copy
  copy.cc
)

target_link_libraries(v2_copy
  evo_test_generated
)

add_executable(v2_validate
  validate.cc
)

target_link_libraries(v2_validate
  evo_test_generated
)
