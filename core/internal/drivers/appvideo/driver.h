#ifndef __CORE_INTERNAL_DRIVERS_GSTREAMER__
#define __CORE_INTERNAL_DRIVERS_GSTREAMER__

#include <pjlib.h>
#include <pjmedia/endpoint.h>
#include <pjmedia/vid_stream.h>
#include <pjmedia-videodev/videodev_imp.h>
#include <pjsua-lib/pjsua.h>

typedef struct vid_dev_factory
{
	pjmedia_vid_dev_factory base;

	pj_pool_t *pool;
	pj_pool_factory *pool_factory;
} vid_dev_factory;

typedef struct vid_dev_stream
{
	pjmedia_vid_dev_stream base;

	pj_pool_t *pool;
	pj_pool_factory *pool_factory;
} vid_dev_stream;

void init_app_video();

extern void go_video_stream_start(void *stream);
extern void go_video_stream_get_frame(pjmedia_vid_dev_stream *stream, pjmedia_frame *frame);
extern void go_video_stream_stop(void *stream);

pj_status_t app_video_dev_factory_init(pjmedia_vid_dev_factory *f);
pj_status_t app_video_dev_factory_destroy(pjmedia_vid_dev_factory *f);
unsigned int app_video_dev_factory_get_dev_count(pjmedia_vid_dev_factory *f);
pj_status_t app_video_dev_factory_get_dev_info(pjmedia_vid_dev_factory *f, unsigned int index, pjmedia_vid_dev_info *info);
pj_status_t app_video_dev_factory_create_stream(pjmedia_vid_dev_factory *f, pjmedia_vid_dev_param *param, const pjmedia_vid_dev_cb *cb, void *user_data, pjmedia_vid_dev_stream **p_vid_strm);
pj_status_t app_video_dev_factory_default_param(pj_pool_t *pool, pjmedia_vid_dev_factory *f, unsigned index, pjmedia_vid_dev_param *param);
pj_status_t app_video_dev_factory_refresh(pjmedia_vid_dev_factory *f);

extern pjmedia_vid_dev_factory_op app_video_dev_factory_op;

pj_status_t app_video_stream_destroy(pjmedia_vid_dev_stream *stream);
pj_status_t app_video_stream_start(pjmedia_vid_dev_stream *stream);
pj_status_t app_video_stream_stop(pjmedia_vid_dev_stream *stream);
pj_status_t app_video_stream_get_cap(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_cap cap, void *value);
pj_status_t app_video_stream_set_cap(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_cap cap, const void *value);
pj_status_t app_video_stream_get_param(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_param *param);
pj_status_t app_video_stream_get_frame(pjmedia_vid_dev_stream *stream, pjmedia_frame *frame);
pj_status_t app_video_stream_put_frame(pjmedia_vid_dev_stream *stream, const pjmedia_frame *frame);

extern pjmedia_vid_dev_stream_op app_video_stream_op;

#endif // __CORE_INTERNAL_DRIVERS_GSTREAMER__