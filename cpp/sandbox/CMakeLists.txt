add_executable(
  sandbox_exec
  sandbox.cc
)

target_link_libraries(
  sandbox_exec
  sandbox_generated
)

add_subdirectory(generated)
