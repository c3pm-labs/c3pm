package raylibadapter

// if ("${CMAKE_SYSTEM_NAME}" MATCHES "Linux")
//     set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -lXxf86vm -lXext -lX11 -lGL -lGLU")
//     find_package(OpenGL REQUIRED)
//     find_package(X11 REQUIRED)
//     find_package(GLUT REQUIRED)
//     find_package(ZLIB REQUIRED)
//     set(LIBRARIES ${RAYLIB_LIBRARY} ${OPENGL_LIBRARIES} ${X11_X11_LIB} ${GLUT_LIBRARY} ${ZLIB_LIBRARIES})
// endif ("${CMAKE_SYSTEM_NAME}" MATCHES "Linux")
var CmakeConfig = `

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
