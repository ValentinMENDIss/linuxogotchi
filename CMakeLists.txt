cmake_minimum_required(VERSION 3.16)

project(Linuxogotchi_Music_Player)

add_executable(${PROJECT_NAME}
  ${CMAKE_CURRENT_SOURCE_DIR}/src/Main.cpp
  ${CMAKE_CURRENT_SOURCE_DIR}/src/AudioDecoder.cpp
  ${CMAKE_CURRENT_SOURCE_DIR}/src/AudioSink.cpp
  ${CMAKE_CURRENT_SOURCE_DIR}/src/MusicPlayer.cpp
)

find_package(ALSA REQUIRED)

find_package(PkgConfig REQUIRED)
pkg_check_modules(MPG123 REQUIRED IMPORTED_TARGET libmpg123)


target_link_libraries(${PROJECT_NAME}
  PRIVATE
    ALSA::ALSA
    PkgConfig::MPG123
)
