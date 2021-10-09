FROM alpine:3.14.2
RUN apk add autoconf automake g++ git libtool make \
 && git clone https://github.com/codership/glb \
 && cd glb/ \
 && ./bootstrap.sh \
 && ./configure \
 && make \
 && make install
# /usr/local/sbin/glbd
# /glb/files/glbd.cfg
#CMD [             \
#  "glbd",         \
#  "--max_conn",   \
#  "1000",         \
#  "--random",     \
#  "--threads",    \
#  "10",           \
#  "--verbose",    \
#  "0.0.0.0:3306"  \
#]