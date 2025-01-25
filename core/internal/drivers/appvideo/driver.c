#include "driver.h"

pjmedia_vid_dev_factory *create_app_video_factory()
{
	pj_pool_factory *pool_factory;
	pj_pool_t *pool;
	vid_dev_factory *factory;

	pool_factory = pjsua_get_pool_factory();
	pool = pj_pool_create(pool_factory, "appvideofc%p", 512, 512, NULL);
	factory = PJ_POOL_ZALLOC_T(pool, vid_dev_factory);

	factory->pool_factory = pjsua_get_pool_factory();
	factory->pool = pool;
	factory->base.op = &app_video_dev_factory_op;

	return &factory->base;
}

void init_app_video()
{
	pjmedia_vid_dev_factory *factory = create_app_video_factory();
	pjmedia_vid_register_factory(NULL, factory);
}

pjmedia_vid_dev_stream_op app_video_stream_op = {
	.destroy = &app_video_stream_destroy,
	.start = &app_video_stream_start,
	.stop = &app_video_stream_stop,
	.get_cap = &app_video_stream_get_cap,
	.set_cap = &app_video_stream_set_cap,
	.get_param = &app_video_stream_get_param,
	.get_frame = &app_video_stream_get_frame,
	.put_frame = &app_video_stream_put_frame,
};

pjmedia_vid_dev_factory_op app_video_dev_factory_op = {
	.init = &app_video_dev_factory_init,
	.destroy = &app_video_dev_factory_destroy,
	.get_dev_count = &app_video_dev_factory_get_dev_count,
	.get_dev_info = &app_video_dev_factory_get_dev_info,
	.create_stream = &app_video_dev_factory_create_stream,
	.default_param = &app_video_dev_factory_default_param,
	.refresh = &app_video_dev_factory_refresh,
};
