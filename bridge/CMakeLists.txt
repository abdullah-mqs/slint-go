cmake_minimum_required(VERSION 3.21)
project(slint_go_bridge LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)

# Slint interpreter ON: .slint file is loaded at runtime (both dev and release use this)
set(SLINT_FEATURE_INTERPRETER ON CACHE BOOL "" FORCE)

# Slint C++ 1.8.0 (stable). Rust/cargo is required for the build (see README).
FetchContent_Declare(
    Slint
    GIT_REPOSITORY https://github.com/slint-ui/slint.git
    GIT_TAG v1.8.0
    GIT_SHALLOW TRUE
    SOURCE_SUBDIR api/cpp
)
FetchContent_MakeAvailable(Slint)

# SHARED lib: All Slint transitive libs are bundled into this DLL/.so,
# so Go only needs -lbridge. A static lib would not link these dependencies.
add_library(bridge SHARED bridge.cpp)
target_include_directories(bridge PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(bridge PUBLIC Slint::Slint)
set_target_properties(bridge PROPERTIES
    POSITION_INDEPENDENT_CODE ON
    WINDOWS_EXPORT_ALL_SYMBOLS ON   # To export C functions on MinGW/Windows
)
if(WIN32)
    target_compile_definitions(bridge PRIVATE NOMINMAX)
endif()
