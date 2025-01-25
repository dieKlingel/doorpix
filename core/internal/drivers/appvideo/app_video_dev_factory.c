#include "driver.h"

pj_status_t app_video_dev_factory_init(pjmedia_vid_dev_factory *f)
{
	return PJ_SUCCESS;
}

unsigned int app_video_dev_factory_get_dev_count(pjmedia_vid_dev_factory *f)
{
	return 1;
}

pj_status_t app_video_dev_factory_get_dev_info(pjmedia_vid_dev_factory *f, unsigned int index, pjmedia_vid_dev_info *info)
{
	pj_ansi_strxcpy(info->name, "DoorPiX Emulated Video Device", 30);
	pj_ansi_strxcpy(info->driver, "DoorPiX", 8);
	info->dir = PJMEDIA_DIR_CAPTURE;
	info->has_callback = PJ_FALSE;
	info->fmt_cnt = 1;
	info->fmt[0].id = PJMEDIA_FORMAT_I420;
	info->fmt[0].detail_type = PJMEDIA_FORMAT_DETAIL_VIDEO;
	info->fmt[0].det.vid.size.w = 640;
	info->fmt[0].det.vid.size.h = 480;
	info->fmt[0].det.vid.fps.num = 30;
	info->fmt[0].det.vid.fps.denum = 1;

	return PJ_SUCCESS;
}

pj_status_t app_video_dev_factory_create_stream(pjmedia_vid_dev_factory *pj_factory, pjmedia_vid_dev_param *param, const pjmedia_vid_dev_cb *cb, void *user_data, pjmedia_vid_dev_stream **p_vid_strm)
{
	vid_dev_factory *factory = (vid_dev_factory *)pj_factory;
	pj_pool_factory *pool_factory = factory->pool_factory;
	pj_pool_t *pool = pj_pool_create(pool_factory, "appvideo%p", 512, 512, NULL);

	vid_dev_stream *stream = PJ_POOL_ZALLOC_T(factory->pool, vid_dev_stream);
	stream->pool_factory = pool_factory;
	stream->pool = pool;
	stream->base.op = &app_video_stream_op;

	// output
	*p_vid_strm = &stream->base;

	return PJ_SUCCESS;
}

pj_status_t app_video_dev_factory_destroy(pjmedia_vid_dev_factory *pj_factory)
{
	vid_dev_factory *factory = (vid_dev_factory *)pj_factory;
	pj_pool_safe_release(&factory->pool);

	return PJ_SUCCESS;
}

pj_status_t app_video_dev_factory_default_param(pj_pool_t *pool, pjmedia_vid_dev_factory *f, unsigned index, pjmedia_vid_dev_param *param)
{
	return PJ_SUCCESS;
}

pj_status_t app_video_dev_factory_refresh(pjmedia_vid_dev_factory *f)
{
	return PJ_SUCCESS;
}
