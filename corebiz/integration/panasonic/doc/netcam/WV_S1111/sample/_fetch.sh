#! /bin/bash

addr='ip:port'
authn='user:pass'

# GETRANDOMNUM - requires access level 1
curl "http://${addr}/cgi-bin/get_randomnum" --digest -u "${authn}" --verbose -o get_randomnum.txt

# info
curl "http://${addr}/cgi-bin/getinfo?FILE=1" --digest -u "${authn}" --verbose -o get_info.txt

# cap
curl "http://${addr}/cgi-bin/get_capability" --digest -u "${authn}" --verbose -o get_cap.txt

# uid - h264
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=h264&reply=info" --digest -u "${authn}" --verbose -o get_uid_h264.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=h264_2&reply=info" --digest -u "${authn}" --verbose -o get_uid_h264_2.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=h264_3&reply=info" --digest -u "${authn}" --verbose -o get_uid_h264_3.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=h264_4&reply=info" --digest -u "${authn}" --verbose -o get_uid_h264_4.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=h264_xxx&reply=info" --digest -u "${authn}" --verbose -o get_uid_h264_xxx.txt

# uid - h265
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=h265&reply=info" --digest -u "${authn}" --verbose -o get_uid_h265.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=h265_2&reply=info" --digest -u "${authn}" --verbose -o get_uid_h265_2.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=h265_3&reply=info" --digest -u "${authn}" --verbose -o get_uid_h265_3.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=h265_4&reply=info" --digest -u "${authn}" --verbose -o get_uid_h265_4.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=h265_xxx&reply=info" --digest -u "${authn}" --verbose -o get_uid_h265_xxx.txt

# uid - jpeg
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=jpeg" --digest -u "${authn}" --verbose -o get_uid_jpeg.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=jpeg_2" --digest -u "${authn}" --verbose -o get_uid_jpeg_2.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=jpeg_3" --digest -u "${authn}" --verbose -o get_uid_jpeg_3.txt
curl "http://${addr}/cgi-bin/getuid?FILE=2&vcodec=jpeg_xxx" --digest -u "${authn}" --verbose -o get_uid_jpeg_xxx.txt

# snapshot - rez
curl "http://${addr}/cgi-bin/camera?resolution=2048" --digest -u "${authn}" --verbose -o snap_rez2048.jpeg
curl "http://${addr}/cgi-bin/camera?resolution=1280" --digest -u "${authn}" --verbose -o snap_rez1280.jpeg
curl "http://${addr}/cgi-bin/camera?resolution=640" --digest -u "${authn}" --verbose -o snap_rez640.jpeg
curl "http://${addr}/cgi-bin/camera?resolution=320" --digest -u "${authn}" --verbose -o snap_rez320.jpeg

# snapshot - stream + ch
curl "http://${addr}/cgi-bin/camera?stream=1&ch=1" --digest -u "${authn}" --verbose -o snap_s1_ch1.jpeg
curl "http://${addr}/cgi-bin/camera?stream=1&ch=2" --digest -u "${authn}" --verbose -o snap_s1_ch2.jpeg
curl "http://${addr}/cgi-bin/camera?stream=1&ch=3" --digest -u "${authn}" --verbose -o snap_s1_ch3.jpeg
curl "http://${addr}/cgi-bin/camera?stream=1&ch=4" --digest -u "${authn}" --verbose -o snap_s1_ch4.jpeg

curl "http://${addr}/cgi-bin/camera?stream=2&ch=1" --digest -u "${authn}" --verbose -o snap_s2_ch1.jpeg
curl "http://${addr}/cgi-bin/camera?stream=3&ch=1" --digest -u "${authn}" --verbose -o snap_s3_ch1.jpeg
curl "http://${addr}/cgi-bin/camera?stream=4&ch=1" --digest -u "${authn}" --verbose -o snap_s4_ch1.jpeg

# sess - requires access level 1
curl "http://${addr}/cgi-bin/man_session?command=get" --digest -u "${authn}" --verbose -o man_session.txt

# keep_alive - this endpoint keeps returning
#    if the param is correct  : 500 Internal server error
#    if the param is incorrect: 403 Forbidden
# curl "http://${addr}/cgi-bin/keep_alive?mode=${stream_mode}&protocol=rtp&UID=${uid}&stream=${stream_num}" --digest -u "${authn}" --verbose
