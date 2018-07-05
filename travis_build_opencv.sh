#!/bin/bash
set -eux -o pipefail

OPENCV_VERSION=${OPENCV_VERSION:-3.4.2}

#GRAPHICAL=ON
GRAPHICAL=${GRAPHICAL:-OFF}

# OpenCV looks for libjpeg in /usr/lib/libjpeg.so, for some reason. However,
# it does not seem to be there in 14.04. Create a link

mkdir -p $HOME/usr/lib

if [[ ! -f "$HOME/usr/lib/libjpeg.so" ]]; then
  ln -s /usr/lib/x86_64-linux-gnu/libjpeg.so $HOME/usr/lib/libjpeg.so
fi

# Same for libpng.so

if [[ ! -f "$HOME/usr/lib/libpng.so" ]]; then
  ln -s /usr/lib/x86_64-linux-gnu/libpng.so $HOME/usr/lib/libpng.so
fi

# Build OpenCV
if [[ ! -e "$HOME/usr/installed-${OPENCV_VERSION}" ]]; then
TMP=$(mktemp -d)
if [[ ! -d "opencv-${OPENCV_VERSION}/build" ]]; then
  curl -sL https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.zip > ${TMP}/opencv.zip
  unzip -q ${TMP}/opencv.zip
  mkdir opencv-${OPENCV_VERSION}/build
  rm ${TMP}/opencv.zip
fi

if [[ ! -d "opencv_contrib-${OPENCV_VERSION}/modules" ]]; then
   curl -sL https://github.com/opencv/opencv_contrib/archive/${OPENCV_VERSION}.zip > ${TMP}/opencv-contrib.zip
   unzip -q ${TMP}/opencv-contrib.zip
   rm ${TMP}/opencv-contrib.zip
fi
rmdir ${TMP}

cd opencv-${OPENCV_VERSION}/build
cmake -D WITH_IPP=${GRAPHICAL} \
      -D WITH_OPENGL=${GRAPHICAL} \
      -D WITH_QT=${GRAPHICAL} \
      -D BUILD_EXAMPLES=OFF \
      -D BUILD_TESTS=OFF \
      -D BUILD_PERF_TESTS=OFF  \
      -D BUILD_opencv_java=OFF \
      -D BUILD_opencv_python=OFF \
      -D BUILD_opencv_python2=OFF \
      -D BUILD_opencv_python3=OFF \
      -D CMAKE_INSTALL_PREFIX=$HOME/usr \
      -D OPENCV_EXTRA_MODULES_PATH=../../opencv_contrib-${OPENCV_VERSION}/modules ..
make -j8
make install && touch $HOME/usr/installed-${OPENCV_VERSION}
cd ../..
touch $HOME/fresh-cache
fi
