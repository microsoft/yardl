add_subdirectory(generated)

add_executable(v1_write
  write.cc
)

target_link_libraries(v1_write
  evo_test_generated_v1
)

add_executable(v1_copy
  copy.cc
)

target_link_libraries(v1_copy
  evo_test_generated_v1
)

add_executable(v1_validate
  validate.cc
)

target_link_libraries(v1_validate
  evo_test_generated_v1
)
