package irrlichtadapter

// set(THREADS_PREFER_PTHREAD_FLAG ON)
// find_package(Threads REQUIRED)
// find_package(X11 REQUIRED)
// find_package(LIBRT rt)
// find_package(OpenGL REQUIRED)
// find_package(GLEW REQUIRED)

// include_directories(${OPENGL_INCLUDE_DIR} ${GLEW_INCLUDE_DIRS})

// set(LIBRARIES ${X11_X11_LIB} ${LIBRT} ${CMAKE_DL_LIBS} m Threads::Threads ${OPENGL_LIBRARIES} ${GLEW_LIBRARIES})

var CmakeConfig = `

if ("${CMAKE_SYSTEM_NAME}" MATCHES "Linux")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -lraylib -lGL -lm -lpthread -ldl -lrt -lX11")
endif ("${CMAKE_SYSTEM_NAME}" MATCHES "Linux")

if (APPLE)
        find_library(CARBON_LIBRARY Carbon)
        find_library(COCOA_LIBRARY Cocoa)
        find_library(IOKIT_LIBRARY IOKit)
        set(OSX_LIBRARIES ${CARBON_LIBRARY} ${COCOA_LIBRARY} ${IOKIT_LIBRARY})
endif (APPLE)

target_include_directories(${C3PM_PROJECT_NAME} PRIVATE src ${IRRLICHT_INCLUDE_DIR})
target_link_libraries(${C3PM_PROJECT_NAME} PUBLIC ${OSX_LIBRARIES} ${LIBRARIES})
`
