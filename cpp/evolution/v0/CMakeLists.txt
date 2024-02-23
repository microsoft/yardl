add_subdirectory(generated)

add_executable(v0_write
  write.cc
)

target_link_libraries(v0_write
  evo_test_generated_v0
)

add_executable(v0_copy
  copy.cc
)

target_link_libraries(v0_copy
  evo_test_generated_v0
)

add_executable(v0_validate
  validate.cc
)

target_link_libraries(v0_validate
  evo_test_generated_v0
)
