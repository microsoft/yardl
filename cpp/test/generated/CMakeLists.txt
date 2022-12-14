# This file was generated by the "yardl" tool. DO NOT EDIT.

# To opt out of generating this file, set cpp.generateCMakeLists to false in the _package.yml file.

# To use the object library defined in this file, add the following to your CMakeLists.txt file:
# target_link_libraries(<your target> test_model_generated)
# add_subdirectory(<path to this directory>)

if(NOT DEFINED CMAKE_CXX_STANDARD OR CMAKE_CXX_STANDARD LESS 20)
  set(USE_DATE true)
  find_package(date REQUIRED)
endif()

if(VCPKG_TARGET_TRIPLET)
  set(HDF5_CXX_LIBRARIES hdf5::hdf5_cpp-shared)
else()
  set(HDF5_CXX_LIBRARIES hdf5::hdf5_cpp)
endif()

find_package(HDF5 REQUIRED COMPONENTS C CXX)
find_package(xtensor REQUIRED)
add_library(test_model_generated OBJECT
  protocols.cc
  binary/protocols.cc
  hdf5/protocols.cc
  mocks.cc
)

target_link_libraries(test_model_generated
  PUBLIC ${HDF5_C_LIBRARIES}
  PUBLIC ${HDF5_CXX_LIBRARIES}
  PUBLIC xtensor
)
if (USE_DATE)
  target_link_libraries(test_model_generated PUBLIC date::date)
endif()
