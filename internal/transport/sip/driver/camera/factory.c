#include "factory.h"

pjmedia_vid_dev_factory_op factory_op = {
	.init = &factory_init,
	.destroy = &factory_destroy,
	.get_dev_count = &factory_get_dev_count,
	.get_dev_info = &factory_get_dev_info,
	.create_stream = &factory_create_stream,
	.default_param = &factory_default_param,
	.refresh = &factory_refresh,
};

pj_status_t factory_init(pjmedia_vid_dev_factory *f)
{
	return PJ_SUCCESS;
}

unsigned int factory_get_dev_count(pjmedia_vid_dev_factory *f)
{
	return 1;
}

pj_status_t factory_get_dev_info(pjmedia_vid_dev_factory *f, unsigned int index, pjmedia_vid_dev_info *info)
{
	factory_t *factory = (factory_t *)f;

	pj_ansi_strxcpy(info->name, factory->options.name, factory->options.name_length + 1);
	pj_ansi_strxcpy(info->driver, factory->options.driver_name, factory->options.driver_name_length + 1);
	info->dir = PJMEDIA_DIR_CAPTURE;
	info->has_callback = PJ_FALSE;
	info->fmt_cnt = 1;
	info->fmt[0].id = PJMEDIA_FORMAT_I420;
	info->fmt[0].detail_type = PJMEDIA_FORMAT_DETAIL_VIDEO;
	info->fmt[0].det.vid.size.w = factory->options.width;
	info->fmt[0].det.vid.size.h = factory->options.height;
	info->fmt[0].det.vid.fps.num = factory->options.framerate;
	info->fmt[0].det.vid.fps.denum = 1;
	info->caps = 0;

	return PJ_SUCCESS;
}

pj_status_t factory_create_stream(pjmedia_vid_dev_factory *pj_factory, pjmedia_vid_dev_param *param, const pjmedia_vid_dev_cb *cb, void *user_data, pjmedia_vid_dev_stream **p_vid_strm)
{
	factory_t *factory = (factory_t *)pj_factory;

	stream_t *stream = PJ_POOL_ZALLOC_T(factory->pool, stream_t);
	stream->pool_factory = factory->pool_factory;
	stream->pool = factory->pool;
	stream->base.op = &stream_op;

	// output
	*p_vid_strm = &stream->base;

	return PJ_SUCCESS;
}

pj_status_t factory_destroy(pjmedia_vid_dev_factory *pj_factory)
{
	factory_t *factory = (factory_t *)pj_factory;
	pj_pool_safe_release(&factory->pool);

	return PJ_SUCCESS;
}

pj_status_t factory_default_param(pj_pool_t *pool, pjmedia_vid_dev_factory *f, unsigned index, pjmedia_vid_dev_param *param)
{
	return PJ_SUCCESS;
}

pj_status_t factory_refresh(pjmedia_vid_dev_factory *f)
{
	return PJ_SUCCESS;
}
