package raylibadapter


// find_package(LIBRT rt)
// find_package(GLEW REQUIRED)
// include_directories(${OPENGL_INCLUDE_DIR} ${GLEW_INCLUDE_DIRS})
var CmakeConfig = `

if ("${CMAKE_SYSTEM_NAME}" MATCHES "Linux")
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -lraylib -lGL -lm -lpthread -ldl -lrt -lX11")
set(THREADS_PREFER_PTHREAD_FLAG ON)
find_package(Threads REQUIRED)
find_package(X11 REQUIRED)
find_package(OpenGL REQUIRED)
set(LIBRARIES ${X11_X11_LIB} ${LIBRT} ${CMAKE_DL_LIBS} m Threads::Threads ${OPENGL_LIBRARIES} ${GLEW_LIBRARIES})
endif ("${CMAKE_SYSTEM_NAME}" MATCHES "Linux")

if (APPLE)
        find_library(CARBON_LIBRARY CoreVideo)
        find_library(COCOA_LIBRARY Cocoa)
        find_library(IOKIT_LIBRARY IOKit)
        find_library(GLUT_LIBRARY GLUT)
        find_library(OPENGL_LIBRARY OpenGL)
        set(OSX_LIBRARIES ${COREVIDEO_LIBRARY} ${COCOA_LIBRARY} ${IOKIT_LIBRARY} ${GLUT_LIBRARY} ${OPENGL_LIBRARY})
        set(MACOSX_DEPLOYMENT_TARGET ${10.9})
endif (APPLE)

target_include_directories(${C3PM_PROJECT_NAME} PRIVITE src ${RAYLIB_INCLUDE_DIR})
target_link_libraries(${C3PM_PROJECT_NAME} PUBLIC ${OSX_LIBRARIES} ${LIBRARIES})
`