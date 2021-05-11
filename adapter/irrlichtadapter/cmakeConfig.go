package irrlichtadapter

var CmakeConfig = `

if ("${CMAKE_SYSTEM_NAME}" MATCHES "Linux")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -lXxf86vm -lXext -lX11 -lGL -lGLU")
endif ("${CMAKE_SYSTEM_NAME}" MATCHES "Linux")
find_package(OpenGL REQUIRED)
find_package(X11 REQUIRED)
find_package(GLUT REQUIRED)
find_package(ZLIB REQUIRED)
set(LIBRARIES ${IRRLICHT_LIBRARY} ${OPENGL_LIBRARIES} ${X11_X11_LIB} ${GLUT_LIBRARY} ${ZLIB_LIBRARIES})

if (APPLE)
        find_library(CARBON_LIBRARY Carbon)
        find_library(COCOA_LIBRARY Cocoa)
        find_library(IOKIT_LIBRARY IOKit)
        set(OSX_LIBRARIES ${CARBON_LIBRARY} ${COCOA_LIBRARY} ${IOKIT_LIBRARY})
endif (APPLE)

target_include_directories(${C3PM_PROJECT_NAME} PRIVATE src ${IRRLICHT_INCLUDE_DIR})
target_link_libraries(${C3PM_PROJECT_NAME} PUBLIC ${OSX_LIBRARIES} ${LIBRARIES})
`
