#include "driver.h"

pj_status_t app_video_stream_get_frame(pjmedia_vid_dev_stream *stream, pjmedia_frame *frame)
{
	go_video_stream_get_frame(stream, frame);
	frame->type = PJMEDIA_FRAME_TYPE_VIDEO;

	return PJ_SUCCESS;
}

pj_status_t app_video_stream_destroy(pjmedia_vid_dev_stream *pj_stream)
{
	vid_dev_stream *stream = (vid_dev_stream *)pj_stream;
	pj_pool_safe_release(&stream->pool);

	return PJ_SUCCESS;
}

pj_status_t app_video_stream_start(pjmedia_vid_dev_stream *stream)
{
	go_video_stream_start(stream);

	return PJ_SUCCESS;
}

pj_status_t app_video_stream_stop(pjmedia_vid_dev_stream *stream)
{
	go_video_stream_stop(stream);

	return PJ_SUCCESS;
}

pj_status_t app_video_stream_get_cap(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_cap cap, void *value)
{
	printf("app_video_stream_get_cap\n");
	return PJ_SUCCESS;
}

pj_status_t app_video_stream_set_cap(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_cap cap, const void *value)
{
	printf("app_video_stream_set_cap\n");
	return PJ_SUCCESS;
}

pj_status_t app_video_stream_get_param(pjmedia_vid_dev_stream *stream, pjmedia_vid_dev_param *param)
{
	printf("app_video_stream_get_param\n");
	return PJ_SUCCESS;
}

pj_status_t app_video_stream_put_frame(pjmedia_vid_dev_stream *stream, const pjmedia_frame *frame)
{
	return PJ_SUCCESS;
}
