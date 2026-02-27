#ifndef DRIVER_H
#define DRIVER_H

#include <pjlib.h>
#include <pjmedia/endpoint.h>
#include <pjmedia/vid_stream.h>
#include <pjmedia-videodev/videodev_imp.h>
#include <pjsua-lib/pjsua.h>

#include "factory.h"

extern int go_stream_start(void *stream);
extern int go_stream_get_frame(pjmedia_vid_dev_stream *stream, pjmedia_frame *frame);
extern int go_stream_stop(void *stream);
int Register(factory_options options);

#endif // DRIVER_H
