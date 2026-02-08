cmake_minimum_required(VERSION 3.16)

project(Linuxogotchi_Music_Player)

add_executable(${PROJECT_NAME}
  ${CMAKE_CURRENT_SOURCE_DIR}/src/main.cpp
)

find_package(ALSA)
find_package(mpg123)

target_link_libraries(${PROJECT_NAME} 
  PRIVATE 
    ALSA::ALSA
    mpg123
)
