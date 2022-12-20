add_executable(
  sandbox_exec
  sandbox.cc
)

target_link_libraries(
  sandbox_exec
  sandbox_generated
)

if(EXISTS ${CMAKE_CURRENT_SOURCE_DIR}/generated)
  add_subdirectory(generated)
endif()
