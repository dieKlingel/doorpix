#ifndef STREAM_H
#define STREAM_H

#include <pjlib.h>
#include <pjmedia/endpoint.h>
#include <pjmedia/vid_stream.h>
#include <pjmedia-videodev/videodev_imp.h>
#include <pjsua-lib/pjsua.h>

pj_status_t stream_destroy(pjmedia_vid_dev_stream *stream);
pj_status_t stream_start(pjmedia_vid_dev_stream *stream);
pj_status_t stream_stop(pjmedia_vid_dev_stream *stream);
pj_status_t stream_get_cap(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_cap cap, void *value);
pj_status_t stream_set_cap(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_cap cap, const void *value);
pj_status_t stream_get_param(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_param *param);
pj_status_t stream_get_frame(pjmedia_vid_dev_stream *stream, pjmedia_frame *frame);
pj_status_t stream_put_frame(pjmedia_vid_dev_stream *stream, const pjmedia_frame *frame);

typedef struct stream
{
	pjmedia_vid_dev_stream base;

	pj_pool_t *pool;
	pj_pool_factory *pool_factory;
} stream_t;

extern pjmedia_vid_dev_stream_op stream_op;

#endif // STREAM_H
